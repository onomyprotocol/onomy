use std::time::Duration;

use common::dockerfile_onomyd;
use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_bank_send, cosmovisor_get_addr, cosmovisor_get_balances, cosmovisor_start,
        set_minimum_gas_price, sh_cosmovisor_no_dbg, wait_for_num_blocks,
    },
    dockerfiles::{dockerfile_hermes, onomy_std_cosmos_daemon_with_arbitrary},
    hermes::{
        hermes_set_gas_price_denom, hermes_start, sh_hermes, write_hermes_config,
        HermesChainConfig, IbcPair,
    },
    onomy_std_init, reprefix_bech32,
    setups::{cosmovisor_add_consumer, marketd_setup, onomyd_setup},
    super_orchestrator::{
        docker::{Container, ContainerNetwork, Dockerfile},
        net_message::NetMessenger,
        remove_files_in_dir, sh,
        stacked_errors::{MapAddError, Result},
        FileOptions, STD_DELAY, STD_TRIES,
    },
    token18, Args, ONOMY_IBC_NOM, TIMEOUT,
};
use tokio::time::sleep;

const CONSUMER_ID: &str = "interchain-security-cd";
const PROVIDER_ACCOUNT_PREFIX: &str = "onomy";
const CONSUMER_ACCOUNT_PREFIX: &str = "cosmos";

#[rustfmt::skip]
const INTERCHAIN_SECURTY_CDD: &str = r#"ENV ICS_VERSION=1.2.0-multiden
ADD https://github.com/cosmos/interchain-security/archive/refs/tags/v$ICS_VERSION.tar.gz /root/v$ICS_VERSION.tar.gz
RUN cd /root && tar -xvf ./v$ICS_VERSION.tar.gz
RUN cd /root/interchain-security-$ICS_VERSION && go build ./cmd/interchain-security-cdd
RUN mkdir -p $DAEMON_HOME/cosmovisor/genesis/$DAEMON_VERSION/bin/
RUN mv /root/interchain-security-$ICS_VERSION/interchain-security-cdd $DAEMON_HOME/cosmovisor/genesis/$DAEMON_VERSION/bin/$DAEMON_NAME
"#;

#[tokio::main]
async fn main() -> Result<()> {
    let args = onomy_std_init()?;

    if let Some(ref s) = args.entry_name {
        match s.as_str() {
            "onomyd" => onomyd_runner(&args).await,
            "consumer" => consumer(&args).await,
            "hermes" => hermes_runner(&args).await,
            _ => format!("entry_name \"{s}\" is not recognized").map_add_err(|| ()),
        }
    } else {
        sh("make build", &[]).await?;
        // copy to dockerfile resources (docker cannot use files from outside cwd)
        sh(
            "cp ./onomyd ./tests/dockerfiles/dockerfile_resources/onomyd",
            &[],
        )
        .await?;
        container_runner(&args).await
    }
}

async fn container_runner(args: &Args) -> Result<()> {
    let logs_dir = "./tests/logs";
    let dockerfiles_dir = "./tests/dockerfiles";
    let bin_entrypoint = &args.bin_name;
    let container_target = "x86_64-unknown-linux-gnu";

    // build internal runner with `--release`
    sh("cargo build --release --bin", &[
        bin_entrypoint,
        "--target",
        container_target,
    ])
    .await?;

    // prepare volumed resources
    remove_files_in_dir("./tests/resources/keyring-test/", &[".address", ".info"]).await?;

    // prepare hermes config
    write_hermes_config(
        &[
            HermesChainConfig::new("onomy", "onomy", false, "anom", true),
            HermesChainConfig::new(CONSUMER_ID, CONSUMER_ACCOUNT_PREFIX, true, "anative", true),
        ],
        &format!("{dockerfiles_dir}/dockerfile_resources"),
    )
    .await?;

    let entrypoint = Some(format!(
        "./target/{container_target}/release/{bin_entrypoint}"
    ));
    let entrypoint = entrypoint.as_deref();

    let mut cn = ContainerNetwork::new(
        "test",
        vec![
            Container::new(
                "hermes",
                Dockerfile::Contents(dockerfile_hermes("__tmp_hermes_config.toml")),
                entrypoint,
                &["--entry-name", "hermes"],
            ),
            Container::new(
                "onomyd",
                Dockerfile::Contents(dockerfile_onomyd()),
                entrypoint,
                &["--entry-name", "onomyd"],
            )
            .volumes(&[(
                "./tests/resources/keyring-test",
                "/root/.onomy/keyring-test",
            )]),
            Container::new(
                "interchain-security-cdd",
                Dockerfile::Contents(onomy_std_cosmos_daemon_with_arbitrary(
                    "interchain-security-cdd",
                    ".interchain-security-cd",
                    "v07-Theta",
                    INTERCHAIN_SECURTY_CDD,
                )),
                entrypoint,
                &["--entry-name", "consumer"],
            )
            .volumes(&[(
                "./tests/resources/keyring-test",
                "/root/.interchain-security-cd/keyring-test",
            )]),
        ],
        Some(dockerfiles_dir),
        true,
        logs_dir,
    )?
    .add_common_volumes(&[(logs_dir, "/logs")]);
    cn.run_all(true).await?;
    cn.wait_with_timeout_all(true, TIMEOUT).await?;
    Ok(())
}

async fn hermes_runner(args: &Args) -> Result<()> {
    let hermes_home = args.hermes_home.as_ref().map_add_err(|| ())?;
    let mut nm_onomyd = NetMessenger::listen_single_connect("0.0.0.0:26000", TIMEOUT).await?;

    // get mnemonic from onomyd
    let mnemonic: String = nm_onomyd.recv().await?;
    // set keys for our chains
    FileOptions::write_str("/root/.hermes/mnemonic.txt", &mnemonic).await?;
    sh_hermes(
        "keys add --chain onomy --mnemonic-file /root/.hermes/mnemonic.txt",
        &[],
    )
    .await?;
    sh_hermes(
        &format!("keys add --chain {CONSUMER_ID} --mnemonic-file /root/.hermes/mnemonic.txt"),
        &[],
    )
    .await?;

    // wait for setup
    nm_onomyd.recv::<()>().await?;

    let ibc_pair = IbcPair::hermes_setup_pair(CONSUMER_ID, "onomy").await?;
    let mut hermes_runner = hermes_start("/logs/hermes_bootstrap_runner.log").await?;
    ibc_pair.hermes_check_acks().await?;

    // tell that chains have been connected
    nm_onomyd.send::<IbcPair>(&ibc_pair).await?;

    // signal to update gas denom
    let ibc_nom = nm_onomyd.recv::<String>().await?;
    hermes_runner.terminate(TIMEOUT).await?;
    hermes_set_gas_price_denom(hermes_home, CONSUMER_ID, &ibc_nom).await?;

    // restart
    let mut hermes_runner = hermes_start("/logs/hermes_runner.log").await?;
    nm_onomyd.send::<()>(&()).await?;

    // termination signal
    nm_onomyd.recv::<()>().await?;
    hermes_runner.terminate(TIMEOUT).await?;
    Ok(())
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let consumer_id = CONSUMER_ID;
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    let mut nm_hermes = NetMessenger::connect(STD_TRIES, STD_DELAY, "hermes:26000")
        .await
        .map_add_err(|| ())?;
    let mut nm_consumer =
        NetMessenger::connect(STD_TRIES, STD_DELAY, &format!("{consumer_id}d:26001"))
            .await
            .map_add_err(|| ())?;

    let mnemonic = onomyd_setup(daemon_home).await?;
    // send mnemonic to hermes
    nm_hermes.send::<String>(&mnemonic).await?;

    // keep these here for local testing purposes
    let addr = &cosmovisor_get_addr("validator").await?;
    sleep(Duration::ZERO).await;

    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", None).await?;

    let ccvconsumer_state = cosmovisor_add_consumer(daemon_home, consumer_id).await?;

    // send to consumer
    nm_consumer.send::<String>(&ccvconsumer_state).await?;

    // send keys
    nm_consumer
        .send::<String>(
            &FileOptions::read_to_string(&format!("{daemon_home}/config/node_key.json")).await?,
        )
        .await?;
    nm_consumer
        .send::<String>(
            &FileOptions::read_to_string(&format!("{daemon_home}/config/priv_validator_key.json"))
                .await?,
        )
        .await?;

    // wait for consumer to be online
    nm_consumer.recv::<()>().await?;
    // notify hermes to connect the chains
    nm_hermes.send::<()>(&()).await?;
    // when hermes is done
    let ibc_pair = nm_hermes.recv::<IbcPair>().await?;
    info!("IbcPair: {ibc_pair:?}");

    // send anom to consumer
    ibc_pair
        .b
        .cosmovisor_ibc_transfer(
            "validator",
            &reprefix_bech32(addr, CONSUMER_ACCOUNT_PREFIX)?,
            &token18(100.0e3, ""),
            "anom",
        )
        .await?;
    // it takes time for the relayer to complete relaying
    wait_for_num_blocks(4).await?;
    // notify consumer that we have sent NOM
    nm_consumer.send::<IbcPair>(&ibc_pair).await?;

    // tell hermes to restart with updated gas denom on its side
    let ibc_nom = nm_consumer.recv::<String>().await?;
    nm_hermes.send::<String>(&ibc_nom).await?;
    nm_hermes.recv::<()>().await?;
    nm_consumer.send::<()>(&()).await?;

    // recieve round trip signal
    nm_consumer.recv::<()>().await?;
    // check that the IBC NOM converted back to regular NOM
    assert_eq!(
        cosmovisor_get_balances("onomy1gk7lg5kd73mcr8xuyw727ys22t7mtz9gh07ul3").await?["anom"],
        "5000"
    );

    // signal to collectively terminate
    nm_hermes.send::<()>(&()).await?;
    nm_consumer.send::<()>(&()).await?;
    cosmovisor_runner.terminate(TIMEOUT).await?;

    FileOptions::write_str(
        "/logs/onomyd_export.json",
        &sh_cosmovisor_no_dbg("export", &[]).await?,
    )
    .await?;

    Ok(())
}

async fn consumer(args: &Args) -> Result<()> {
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    let chain_id = CONSUMER_ID;
    let mut nm_onomyd = NetMessenger::listen_single_connect("0.0.0.0:26001", TIMEOUT).await?;
    // we need the initial consumer state
    let ccvconsumer_state_s: String = nm_onomyd.recv().await?;

    marketd_setup(daemon_home, chain_id, &ccvconsumer_state_s).await?;
    // make sure switching is possible
    set_minimum_gas_price(daemon_home, "1anative").await?;

    // get keys
    let node_key = nm_onomyd.recv::<String>().await?;
    // we used same keys for consumer as producer, need to copy them over or else
    // the node will not be a working validator for itself
    FileOptions::write_str(&format!("{daemon_home}/config/node_key.json"), &node_key).await?;

    let priv_validator_key = nm_onomyd.recv::<String>().await?;
    FileOptions::write_str(
        &format!("{daemon_home}/config/priv_validator_key.json"),
        &priv_validator_key,
    )
    .await?;

    let mut cosmovisor_runner =
        cosmovisor_start(&format!("{chain_id}d_bootstrap_runner.log"), None).await?;

    let addr = &cosmovisor_get_addr("validator").await?;

    // signal that we have started
    nm_onomyd.send::<()>(&()).await?;

    // wait for producer to send us stuff
    let ibc_pair = nm_onomyd.recv::<IbcPair>().await?;
    // get the name of the IBC NOM. Note that we can't do this on the onomyd side,
    // it has to be with respect to the consumer side
    let ibc_nom = &ibc_pair.a.get_ibc_denom("anom").await?;
    assert_eq!(ibc_nom, ONOMY_IBC_NOM);
    let balances = cosmovisor_get_balances(addr).await?;
    assert!(balances.contains_key(ibc_nom));

    // we have IBC NOM, shut down, change gas in app.toml, restart
    cosmovisor_runner.terminate(TIMEOUT).await?;
    set_minimum_gas_price(daemon_home, &format!("1{ibc_nom}")).await?;
    let mut cosmovisor_runner = cosmovisor_start(&format!("{chain_id}d_runner.log"), None).await?;
    // tell hermes to restart with updated gas denom on its side
    nm_onomyd.send::<String>(ibc_nom).await?;
    nm_onomyd.recv::<()>().await?;
    info!("restarted with new gas denom");

    // test normal transfer
    let dst_addr = &reprefix_bech32(
        "onomy1gk7lg5kd73mcr8xuyw727ys22t7mtz9gh07ul3",
        CONSUMER_ACCOUNT_PREFIX,
    )?;
    cosmovisor_bank_send(addr, dst_addr, "5000", ibc_nom).await?;
    assert_eq!(cosmovisor_get_balances(dst_addr).await?[ibc_nom], "5000");

    let test_addr = &reprefix_bech32(
        "onomy1gk7lg5kd73mcr8xuyw727ys22t7mtz9gh07ul3",
        PROVIDER_ACCOUNT_PREFIX,
    )?;
    info!("sending back to {}", test_addr);

    // send some IBC NOM back to origin chain using it as gas
    ibc_pair
        .a
        .cosmovisor_ibc_transfer("validator", test_addr, "5000", ibc_nom)
        .await?;
    wait_for_num_blocks(4).await?;

    // round trip signal
    nm_onomyd.send::<()>(&()).await?;

    // termination signal
    nm_onomyd.recv::<()>().await?;
    cosmovisor_runner.terminate(TIMEOUT).await?;

    FileOptions::write_str(
        &format!("/logs/{chain_id}_export.json"),
        &sh_cosmovisor_no_dbg("export", &[]).await?,
    )
    .await?;

    Ok(())
}
