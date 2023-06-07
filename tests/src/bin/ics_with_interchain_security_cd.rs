use std::time::Duration;

use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_get_addr, cosmovisor_start, fast_block_times, onomyd_setup, sh_cosmovisor,
        wait_for_height,
    },
    hermes::{create_channel_pair, create_connection_pair, sh_hermes},
    json_inner, onomy_std_init, reprefix_bech32,
    super_orchestrator::{
        docker::{Container, ContainerNetwork},
        net_message::NetMessenger,
        remove_files_in_dir, sh,
        stacked_errors::{MapAddError, Result},
        Command, FileOptions, STD_DELAY, STD_TRIES,
    },
    token18, Args, TIMEOUT,
};
use serde_json::Value;
use tokio::time::sleep;

#[tokio::main]
async fn main() -> Result<()> {
    let args = onomy_std_init()?;

    if let Some(ref s) = args.entry_name {
        match s.as_str() {
            "onomyd" => onomyd_runner(&args).await,
            "interchain-security-cd" => interchain_security_cd_runner(&args).await,
            "hermes" => hermes_runner().await,
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
    let bin_entrypoint = &args.bin_name;
    let container_target = "x86_64-unknown-linux-gnu";
    let logs_dir = "./tests/logs";

    // build internal runner with `--release`
    sh("cargo build --release --bin", &[
        bin_entrypoint,
        "--target",
        container_target,
    ])
    .await?;

    // prepare volumed resources
    remove_files_in_dir("./tests/resources/keyring-test/", &["address", "info"]).await?;

    let entrypoint = Some(format!(
        "./target/{container_target}/release/{bin_entrypoint}"
    ));
    let entrypoint = entrypoint.as_deref();
    let volumes = vec![(logs_dir, "/logs")];
    let mut onomyd_volumes = volumes.clone();
    let mut consumer_volumes = volumes.clone();
    onomyd_volumes.push((
        "./tests/resources/keyring-test",
        "/root/.onomy/keyring-test",
    ));
    consumer_volumes.push((
        "./tests/resources/keyring-test",
        "/root/.interchain-security-c/keyring-test",
    ));

    let mut cn = ContainerNetwork::new(
        "test",
        vec![
            Container::new(
                "hermes",
                Some("./tests/dockerfiles/hermes.dockerfile"),
                None,
                &volumes,
                entrypoint,
                &["--entry-name", "hermes"],
            ),
            Container::new(
                "onomyd",
                Some("./tests/dockerfiles/onomyd.dockerfile"),
                None,
                &onomyd_volumes,
                entrypoint,
                &["--entry-name", "onomyd"],
            ),
            Container::new(
                "interchain-security-cd",
                Some("./tests/dockerfiles/interchain_security_cd.dockerfile"),
                None,
                &consumer_volumes,
                entrypoint,
                &["--entry-name", "interchain-security-cd"],
            ),
        ],
        true,
        logs_dir,
    )?;
    cn.run_all(true).await?;
    cn.wait_with_timeout_all(true, TIMEOUT).await?;
    Ok(())
}

async fn hermes_runner() -> Result<()> {
    let mut nm_onomyd = NetMessenger::listen_single_connect("0.0.0.0:26000", TIMEOUT).await?;

    let mnemonic: String = nm_onomyd.recv().await?;
    // set keys for our chains
    FileOptions::write_str("/root/.hermes/mnemonic.txt", &mnemonic).await?;
    sh_hermes(
        "keys add --chain onomy --mnemonic-file /root/.hermes/mnemonic.txt",
        &[],
    )
    .await?;
    sh_hermes(
        "keys add --chain interchain-security-c --mnemonic-file /root/.hermes/mnemonic.txt",
        &[],
    )
    .await?;

    nm_onomyd.recv::<()>().await?;

    // https://hermes.informal.systems/tutorials/local-chains/add-a-new-relay-path.html

    // Note: For ICS, there is a point where a handshake must be initiated by the
    // consumer chain, so we must make the consumer chain the "a-chain" and the
    // producer chain the "b-chain"

    let b_chain = "onomy";
    let a_chain = "interchain-security-c";
    // a client is already created because of the ICS setup
    //let client_pair = create_client_pair(a_chain, b_chain).await?;
    // create one client and connection pair that will be used for IBC transfer and
    // ICS communication
    let connection_pair = create_connection_pair(a_chain, b_chain).await?;

    // a_chain<->b_chain transfer<->transfer
    let transfer_channel_pair =
        create_channel_pair(a_chain, &connection_pair.0, "transfer", "transfer", false).await?;

    // a_chain<->b_chain consumer<->provider
    let ics_channel_pair =
        create_channel_pair(a_chain, &connection_pair.0, "consumer", "provider", true).await?;

    let hermes_log = FileOptions::write2("/logs", "hermes_runner.log");
    let mut hermes_runner = Command::new("hermes start", &[])
        .stderr_log(&hermes_log)
        .stdout_log(&hermes_log)
        .run()
        .await?;

    info!("Onomy Network has been setup");

    sleep(Duration::from_secs(5)).await;

    // check all channels on both sides
    sh_hermes("query packet acks --chain", &[
        b_chain,
        "--port",
        "transfer",
        "--channel",
        &transfer_channel_pair.0,
    ])
    .await?;
    sh_hermes("query packet acks --chain", &[
        a_chain,
        "--port",
        "transfer",
        "--channel",
        &transfer_channel_pair.1,
    ])
    .await?;
    sh_hermes("query packet acks --chain", &[
        b_chain,
        "--port",
        "provider",
        "--channel",
        &ics_channel_pair.0,
    ])
    .await?;
    sh_hermes("query packet acks --chain", &[
        a_chain,
        "--port",
        "consumer",
        "--channel",
        &ics_channel_pair.1,
    ])
    .await?;

    //hermes tx ft-transfer --timeout-seconds 10 --dst-chain interchain-security-c
    // --src-chain onomy --src-port transfer --src-channel channel-0 --amount
    // 1337 --denom anom

    nm_onomyd.send::<()>(&()).await?;

    hermes_runner.terminate().await?;
    Ok(())
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let consumer_id = "interchain-security-c";
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    let mut nm_hermes = NetMessenger::connect(STD_TRIES, STD_DELAY, "hermes:26000")
        .await
        .map_add_err(|| ())?;
    let mut nm_consumer =
        NetMessenger::connect(STD_TRIES, STD_DELAY, "interchain-security-cd:26001")
            .await
            .map_add_err(|| ())?;

    let mnemonic = onomyd_setup(daemon_home, false).await?;

    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", true, None).await?;

    let proposal_id = "1";

    // TODO we think we will make the redistribution fraction 0 and either make a
    // native "native" or IBC NOM as the gas denom (may take a gov proposal for
    // bootstrap)

    // `json!` doesn't like large literals beyond i32
    let proposal_s = &format!(
        r#"{{
        "title": "Propose the addition of a new chain",
        "description": "add consumer chain",
        "chain_id": "{consumer_id}",
        "initial_height": {{
            "revision_number": 0,
            "revision_height": 1
        }},
        "genesis_hash": "Z2VuX2hhc2g=",
        "binary_hash": "YmluX2hhc2g=",
        "spawn_time": "2023-05-18T01:15:49.83019476-05:00",
        "consumer_redistribution_fraction": "0.0",
        "blocks_per_distribution_transmission": 1000,
        "historical_entries": 10000,
        "ccv_timeout_period": 2419200000000000,
        "transfer_timeout_period": 3600000000000,
        "unbonding_period": 1728000000000000,
        "deposit": "2000000000000000000000anom"
    }}"#
    );
    // we will just place the file under the config folder
    let proposal_file_path = format!("{daemon_home}/config/consumer_add_proposal.json");
    FileOptions::write_str(&proposal_file_path, proposal_s)
        .await
        .map_add_err(|| ())?;

    let gas_args = [
        "--gas",
        "auto",
        "--gas-adjustment",
        "1.3",
        "-y",
        "-b",
        "block",
        "--from",
        "validator",
    ]
    .as_slice();
    sh_cosmovisor(
        "tx gov submit-proposal consumer-addition",
        &[&[proposal_file_path.as_str()], gas_args].concat(),
    )
    .await?;
    // the deposit is done as part of the chain addition proposal
    sh_cosmovisor(
        "tx gov vote",
        &[[proposal_id, "yes"].as_slice(), gas_args].concat(),
    )
    .await?;

    // In the mean time get consensus key assignment done

    let tendermint_key: Value = serde_json::from_str(
        &FileOptions::read_to_string(&format!("{daemon_home}/config/priv_validator_key.json"))
            .await?,
    )?;
    let tendermint_key = json_inner(&tendermint_key["pub_key"]["value"]);
    let tendermint_key =
        format!("{{\"@type\":\"/cosmos.crypto.ed25519.PubKey\",\"key\":\"{tendermint_key}\"}}");

    // do this before getting the consumer-genesis
    sh_cosmovisor(
        "tx provider assign-consensus-key",
        &[[consumer_id, tendermint_key.as_str()].as_slice(), gas_args].concat(),
    )
    .await?;

    wait_for_height(STD_TRIES, STD_DELAY, 5).await?;

    let ccvconsumer_state = sh_cosmovisor("query provider consumer-genesis", &[
        consumer_id,
        "-o",
        "json",
    ])
    .await?;

    nm_hermes.send::<String>(&mnemonic).await?;

    // send to consumer
    nm_consumer.send::<String>(&ccvconsumer_state).await?;

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
    nm_hermes.recv::<()>().await?;
    // finish
    nm_consumer.send::<()>(&()).await?;

    //cosmovisor("tx ibc-transfer transfer", &[port, channel, receiver,
    // amount]).await?;

    cosmovisor_runner.terminate().await?;
    Ok(())
}

async fn interchain_security_cd_runner(args: &Args) -> Result<()> {
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    let mut nm_onomyd = NetMessenger::listen_single_connect("0.0.0.0:26001", TIMEOUT).await?;
    let chain_id = "interchain-security-c";
    sh_cosmovisor("config chain-id", &[chain_id]).await?;
    sh_cosmovisor("config keyring-backend test", &[]).await?;
    sh_cosmovisor("init --overwrite", &[chain_id]).await?;
    let genesis_file_path = format!("{daemon_home}/config/genesis.json");

    // we need the initial consumer state
    let ccvconsumer_state_s: String = nm_onomyd.recv().await?;
    let ccvconsumer_state: Value = serde_json::from_str(&ccvconsumer_state_s)?;

    // add `ccvconsumer_state` to genesis

    let genesis_s = FileOptions::read_to_string(&genesis_file_path).await?;

    let mut genesis: Value = serde_json::from_str(&genesis_s)?;
    genesis["app_state"]["ccvconsumer"] = ccvconsumer_state;
    let genesis_s = genesis.to_string();

    // I will name the token "native" because it won't be staked in the normal sense
    let genesis_s = genesis_s.replace("\"stake\"", "\"native\"");

    FileOptions::write_str(&genesis_file_path, &genesis_s).await?;

    let addr: &String = &cosmovisor_get_addr("validator").await?;
    // use "cosmos" prefix on the consumer chain so that we don't have to modify it
    let addr = &reprefix_bech32(addr, "cosmos").unwrap();

    // we need some native token in the bank, and don't need gentx
    sh_cosmovisor("add-genesis-account", &[addr, &token18(2.0e6, "native")]).await?;

    FileOptions::write_str(
        &format!("/logs/{chain_id}_genesis.json"),
        &FileOptions::read_to_string(&genesis_file_path).await?,
    )
    .await?;

    fast_block_times(daemon_home).await?;

    // we used same keys for consumer as producer, need to copy them over or else
    // the node will not be a working validator for itself
    FileOptions::write_str(
        &format!("{daemon_home}/config/node_key.json"),
        &nm_onomyd.recv::<String>().await?,
    )
    .await?;
    FileOptions::write_str(
        &format!("{daemon_home}/config/priv_validator_key.json"),
        &nm_onomyd.recv::<String>().await?,
    )
    .await?;

    let mut cosmovisor_runner =
        cosmovisor_start(&format!("{chain_id}d_runner.log"), true, None).await?;

    // signal that we have started
    nm_onomyd.send::<()>(&()).await?;

    nm_onomyd.recv::<()>().await?;

    cosmovisor_runner.terminate().await?;
    Ok(())
}
