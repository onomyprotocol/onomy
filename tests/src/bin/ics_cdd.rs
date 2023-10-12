use std::time::Duration;

use common::dockerfile_onomyd;
use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_bank_send, cosmovisor_get_addr, cosmovisor_get_balances, cosmovisor_start,
        set_minimum_gas_price, sh_cosmovisor, sh_cosmovisor_no_dbg, sh_cosmovisor_tx,
        wait_for_num_blocks,
    },
    dockerfiles::{dockerfile_hermes, onomy_std_cosmos_daemon_with_arbitrary},
    hermes::{
        hermes_set_gas_price_denom, hermes_start, sh_hermes, write_hermes_config,
        HermesChainConfig, IbcPair,
    },
    onomy_std_init, reprefix_bech32,
    setups::{
        cosmovisor_add_consumer, marketd_setup, onomyd_setup, test_proposal, CosmosSetupOptions,
    },
    super_orchestrator::{
        docker::{Container, ContainerNetwork, Dockerfile},
        net_message::NetMessenger,
        remove_files_in_dir, sh,
        stacked_errors::{Error, Result, StackableErr},
        FileOptions, STD_DELAY, STD_TRIES,
    },
    token18, u64_array_bigints,
    u64_array_bigints::u256,
    Args, ONOMY_IBC_NOM, TIMEOUT,
};
use tokio::time::sleep;

const CONSUMER_ID: &str = "interchain-security-cd";
const CONSUMER_VERSION: &str = "v07-Theta";
const PROVIDER_ACCOUNT_PREFIX: &str = "onomy";
const CONSUMER_ACCOUNT_PREFIX: &str = "cosmos";

fn consumer_binary_name() -> String {
    format!("{CONSUMER_ID}d")
}

fn consumer_directory() -> String {
    format!(".{CONSUMER_ID}")
    //format!(".onomy_{CONSUMER_ID}")
}

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
            _ => Err(Error::from(format!("entry_name \"{s}\" is not recognized"))),
        }
    } else {
        sh("make build", &[]).await.stack()?;
        // copy to dockerfile resources (docker cannot use files from outside cwd)
        sh(
            "cp ./onomyd ./tests/dockerfiles/dockerfile_resources/onomyd",
            &[],
        )
        .await
        .stack()?;
        container_runner(&args).await.stack()
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
    .await
    .stack()?;

    // prepare volumed resources
    remove_files_in_dir("./tests/resources/keyring-test/", &[".address", ".info"])
        .await
        .stack()?;

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
                &consumer_binary_name(),
                Dockerfile::Contents(onomy_std_cosmos_daemon_with_arbitrary(
                    &consumer_binary_name(),
                    &consumer_directory(),
                    CONSUMER_VERSION,
                    INTERCHAIN_SECURTY_CDD,
                )),
                entrypoint,
                &["--entry-name", "consumer"],
            )
            .volumes(&[(
                "./tests/resources/keyring-test",
                &format!("/root/{}/keyring-test", consumer_directory()),
            )]),
        ],
        Some(dockerfiles_dir),
        true,
        logs_dir,
    )
    .stack()?;
    cn.add_common_volumes(&[(logs_dir, "/logs")]);
    let uuid = cn.uuid_as_string();
    cn.add_common_entrypoint_args(&["--uuid", &uuid]);

    // prepare hermes config
    write_hermes_config(
        &[
            HermesChainConfig::new(
                "onomy",
                &format!("onomyd_{uuid}"),
                "onomy",
                false,
                "anom",
                true,
            ),
            HermesChainConfig::new(
                CONSUMER_ID,
                &format!("{}_{}", consumer_binary_name(), uuid),
                CONSUMER_ACCOUNT_PREFIX,
                true,
                "anative",
                true,
            ),
        ],
        &format!("{dockerfiles_dir}/dockerfile_resources"),
    )
    .await
    .stack()?;

    cn.run_all(true).await.stack()?;
    cn.wait_with_timeout_all(true, TIMEOUT).await.stack()?;
    cn.terminate_all().await;
    Ok(())
}

async fn hermes_runner(args: &Args) -> Result<()> {
    let hermes_home = args.hermes_home.as_ref().stack()?;
    let mut nm_onomyd = NetMessenger::listen_single_connect("0.0.0.0:26000", TIMEOUT)
        .await
        .stack()?;

    // get mnemonic from onomyd
    let mnemonic: String = nm_onomyd.recv().await.stack()?;
    // set keys for our chains
    FileOptions::write_str("/root/.hermes/mnemonic.txt", &mnemonic)
        .await
        .stack()?;
    sh_hermes(
        "keys add --chain onomy --mnemonic-file /root/.hermes/mnemonic.txt",
        &[],
    )
    .await
    .stack()?;
    sh_hermes(
        &format!("keys add --chain {CONSUMER_ID} --mnemonic-file /root/.hermes/mnemonic.txt"),
        &[],
    )
    .await
    .stack()?;

    // wait for setup
    nm_onomyd.recv::<()>().await.stack()?;

    let ibc_pair = IbcPair::hermes_setup_ics_pair(CONSUMER_ID, "onomy")
        .await
        .stack()?;
    let mut hermes_runner = hermes_start("/logs/hermes_bootstrap_runner.log")
        .await
        .stack()?;
    ibc_pair.hermes_check_acks().await.stack()?;

    // tell that chains have been connected
    nm_onomyd.send::<IbcPair>(&ibc_pair).await.stack()?;

    // signal to update gas denom
    let ibc_nom = nm_onomyd.recv::<String>().await.stack()?;
    hermes_runner.terminate(TIMEOUT).await.stack()?;
    hermes_set_gas_price_denom(hermes_home, CONSUMER_ID, &ibc_nom)
        .await
        .stack()?;

    // restart
    let mut hermes_runner = hermes_start("/logs/hermes_runner.log").await.stack()?;
    nm_onomyd.send::<()>(&()).await.stack()?;

    // termination signal
    nm_onomyd.recv::<()>().await.stack()?;
    hermes_runner.terminate(TIMEOUT).await.stack()?;
    Ok(())
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let consumer_id = CONSUMER_ID;
    let daemon_home = args.daemon_home.as_ref().stack()?;
    let mut nm_hermes =
        NetMessenger::connect(STD_TRIES, STD_DELAY, &format!("hermes_{}:26000", args.uuid))
            .await
            .stack()?;
    let mut nm_consumer = NetMessenger::connect(
        STD_TRIES,
        STD_DELAY,
        &format!("{}_{}:26001", consumer_binary_name(), args.uuid),
    )
    .await
    .stack()
    .stack()?;

    let mnemonic = onomyd_setup(CosmosSetupOptions::new(daemon_home))
        .await
        .stack()?;
    // send mnemonic to hermes
    nm_hermes.send::<String>(&mnemonic).await.stack()?;

    // keep these here for local testing purposes
    let addr = &cosmovisor_get_addr("validator").await.stack()?;
    sleep(Duration::ZERO).await;

    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", None).await.stack()?;

    let ccvconsumer_state = cosmovisor_add_consumer(
        daemon_home,
        consumer_id,
        &test_proposal(consumer_id, "anative"),
    )
    .await
    .stack()?;

    // send to consumer
    nm_consumer
        .send::<String>(&ccvconsumer_state)
        .await
        .stack()?;

    // send keys
    nm_consumer
        .send::<String>(
            &FileOptions::read_to_string(&format!("{daemon_home}/config/node_key.json"))
                .await
                .stack()?,
        )
        .await
        .stack()?;
    nm_consumer
        .send::<String>(
            &FileOptions::read_to_string(&format!("{daemon_home}/config/priv_validator_key.json"))
                .await
                .stack()?,
        )
        .await
        .stack()?;

    // wait for consumer to be online
    nm_consumer.recv::<()>().await.stack()?;
    // notify hermes to connect the chains
    nm_hermes.send::<()>(&()).await.stack()?;
    // when hermes is done
    let ibc_pair = nm_hermes.recv::<IbcPair>().await.stack()?;
    info!("IbcPair: {ibc_pair:?}");

    // send anom to consumer
    ibc_pair
        .b
        .cosmovisor_ibc_transfer(
            "validator",
            &reprefix_bech32(addr, CONSUMER_ACCOUNT_PREFIX).stack()?,
            &token18(1.0e3, ""),
            "anom",
        )
        .await
        .stack()?;
    // it takes time for the relayer to complete relaying
    wait_for_num_blocks(5).await.stack()?;
    // notify consumer that we have sent NOM
    nm_consumer.send::<IbcPair>(&ibc_pair).await.stack()?;

    // tell hermes to restart with updated gas denom on its side
    let ibc_nom = nm_consumer.recv::<String>().await.stack()?;
    nm_hermes.send::<String>(&ibc_nom).await.stack()?;
    nm_hermes.recv::<()>().await.stack()?;
    nm_consumer.send::<()>(&()).await.stack()?;

    // recieve round trip signal
    nm_consumer.recv::<()>().await.stack()?;
    // check that the IBC NOM converted back to regular NOM
    assert_eq!(
        cosmovisor_get_balances("onomy1gk7lg5kd73mcr8xuyw727ys22t7mtz9gh07ul3")
            .await
            .stack()?["anom"],
        u256!(5000)
    );

    // signal to collectively terminate
    nm_hermes.send::<()>(&()).await.stack()?;
    nm_consumer.send::<()>(&()).await.stack()?;
    cosmovisor_runner.terminate(TIMEOUT).await.stack()?;

    FileOptions::write_str(
        "/logs/onomyd_export.json",
        &sh_cosmovisor_no_dbg("export", &[]).await.stack()?,
    )
    .await
    .stack()?;

    Ok(())
}

async fn consumer(args: &Args) -> Result<()> {
    let daemon_home = args.daemon_home.as_ref().stack()?;
    let chain_id = CONSUMER_ID;
    let mut nm_onomyd = NetMessenger::listen_single_connect("0.0.0.0:26001", TIMEOUT)
        .await
        .stack()?;
    // we need the initial consumer state
    let ccvconsumer_state_s: String = nm_onomyd.recv().await.stack()?;

    marketd_setup(daemon_home, chain_id, &ccvconsumer_state_s)
        .await
        .stack()?;

    // get keys
    let node_key = nm_onomyd.recv::<String>().await.stack()?;
    // we used same keys for consumer as provider, need to copy them over or else
    // the node will not be a working validator for itself
    FileOptions::write_str(&format!("{daemon_home}/config/node_key.json"), &node_key)
        .await
        .stack()?;

    let priv_validator_key = nm_onomyd.recv::<String>().await.stack()?;
    FileOptions::write_str(
        &format!("{daemon_home}/config/priv_validator_key.json"),
        &priv_validator_key,
    )
    .await
    .stack()?;

    let mut cosmovisor_runner =
        cosmovisor_start(&format!("{chain_id}d_bootstrap_runner.log"), None)
            .await
            .stack()?;

    let addr = &cosmovisor_get_addr("validator").await.stack()?;

    // signal that we have started
    nm_onomyd.send::<()>(&()).await.stack()?;

    // wait for provider to send us stuff
    let ibc_pair = nm_onomyd.recv::<IbcPair>().await.stack()?;
    // get the name of the IBC NOM. Note that we can't do this on the onomyd side,
    // it has to be with respect to the consumer side
    let ibc_nom = &ibc_pair.a.get_ibc_denom("anom").await.stack()?;
    assert_eq!(ibc_nom, ONOMY_IBC_NOM);
    let balances = cosmovisor_get_balances(addr).await.stack()?;
    assert!(balances.contains_key(ibc_nom));

    // we have IBC NOM, shut down, change gas in app.toml, restart
    cosmovisor_runner.terminate(TIMEOUT).await.stack()?;
    set_minimum_gas_price(daemon_home, &format!("1{ibc_nom}"))
        .await
        .stack()?;
    let mut cosmovisor_runner = cosmovisor_start(&format!("{chain_id}d_runner.log"), None)
        .await
        .stack()?;
    // tell hermes to restart with updated gas denom on its side
    nm_onomyd.send::<String>(ibc_nom).await.stack()?;
    nm_onomyd.recv::<()>().await.stack()?;
    info!("restarted with new gas denom");

    // test normal transfer
    let dst_addr = &reprefix_bech32(
        "onomy1gk7lg5kd73mcr8xuyw727ys22t7mtz9gh07ul3",
        CONSUMER_ACCOUNT_PREFIX,
    )
    .stack()?;
    cosmovisor_bank_send(addr, dst_addr, "5000", ibc_nom)
        .await
        .stack()?;
    assert_eq!(
        cosmovisor_get_balances(dst_addr).await.stack()?[ibc_nom],
        u256!(5000)
    );

    let test_addr = &reprefix_bech32(
        "onomy1gk7lg5kd73mcr8xuyw727ys22t7mtz9gh07ul3",
        PROVIDER_ACCOUNT_PREFIX,
    )
    .stack()?;
    info!("sending back to {}", test_addr);

    // avoid conflict with hermes relayer
    wait_for_num_blocks(5).await.stack()?;

    // send some IBC NOM back to origin chain using it as gas
    ibc_pair
        .a
        .cosmovisor_ibc_transfer("validator", test_addr, "5000", ibc_nom)
        .await
        .stack()?;
    wait_for_num_blocks(5).await.stack()?;

    let pubkey = sh_cosmovisor("tendermint show-validator", &[])
        .await
        .stack()?;
    let pubkey = pubkey.trim();
    sh_cosmovisor_tx("staking", &[
        "create-validator",
        "--commission-max-change-rate",
        "0.01",
        "--commission-max-rate",
        "0.10",
        "--commission-rate",
        "0.05",
        "--from",
        "validator",
        "--min-self-delegation",
        "1",
        "--amount",
        &token18(500.0, ONOMY_IBC_NOM),
        "--fees",
        &format!("1000000{ONOMY_IBC_NOM}"),
        "--pubkey",
        pubkey,
        "-y",
        "-b",
        "block",
    ])
    .await
    .stack()?;

    // round trip signal
    nm_onomyd.send::<()>(&()).await.stack()?;

    // termination signal
    nm_onomyd.recv::<()>().await.stack()?;

    // interchain-security-cd does not support this proposal
    /*
    // but first, test governance with IBC NOM as the token
    let test_crisis_denom = ONOMY_IBC_NOM;
    let test_deposit = token18(2000.0, ONOMY_IBC_NOM);
    wait_for_num_blocks(1).await.stack()?;
    cosmovisor_gov_file_proposal(
        daemon_home,
        "param-change",
        &format!(
            r#"
    {{
        "title": "Parameter Change",
        "description": "Making a parameter change",
        "changes": [
          {{
            "subspace": "crisis",
            "key": "ConstantFee",
            "value": {{"denom":"{test_crisis_denom}","amount":"1337"}}
          }}
        ],
        "deposit": "{test_deposit}"
    }}
    "#
        ),
        &format!("1{ibc_nom}"),
    )
    .await
    .stack()?;
    wait_for_num_blocks(5).await.stack()?;
    // just running this for debug, param querying is weird because it is json
    // inside of yaml, so we will instead test the exported genesis
    sh_cosmovisor("query params subspace crisis ConstantFee", &[])
        .await
        .stack()?;*/

    cosmovisor_runner.terminate(TIMEOUT).await.stack()?;

    let exported = sh_cosmovisor_no_dbg("export", &[]).await.stack()?;
    FileOptions::write_str(&format!("/logs/{chain_id}_export.json"), &exported)
        .await
        .stack()?;
    /*let exported = yaml_str_to_json_value(&exported).stack()?;
    assert_eq!(
        exported["app_state"]["crisis"]["constant_fee"]["denom"],
        test_crisis_denom
    );
    assert_eq!(
        exported["app_state"]["crisis"]["constant_fee"]["amount"],
        "1337"
    );*/

    Ok(())
}
