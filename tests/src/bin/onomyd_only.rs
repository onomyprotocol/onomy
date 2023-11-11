//! based from onomy_tests/tests/src/bin/onomyd_only.r

use std::time::Duration;

use common::{container_runner, dockerfile_onomyd};
use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_get_addr, cosmovisor_gov_file_proposal, cosmovisor_start, get_apr_annual,
        get_delegations_to, get_staking_pool, get_treasury, get_treasury_inflation_annual,
        sh_cosmovisor, sh_cosmovisor_no_debug, sh_cosmovisor_tx, wait_for_num_blocks,
    },
    onomy_std_init, reprefix_bech32,
    setups::{onomyd_setup, CosmosSetupOptions},
    super_orchestrator::{
        sh,
        stacked_errors::{ensure, ensure_eq, Error, Result, StackableErr},
        stacked_get, FileOptions,
    },
    token18, yaml_str_to_json_value, Args, ONOMY_IBC_NOM, TIMEOUT,
};
use serde_json::json;
use tokio::time::sleep;

#[tokio::main]
async fn main() -> Result<()> {
    let args = onomy_std_init()?;

    if let Some(ref s) = args.entry_name {
        match s.as_str() {
            "onomyd" => onomyd_runner(&args).await,
            _ => Err(Error::from(format!("entry_name \"{s}\" is not recognized"))),
        }
    } else {
        sh(["make build"]).await.stack()?;
        // copy to dockerfile resources (docker cannot use files from outside cwd)
        sh(["cp ./onomyd ./tests/dockerfiles/dockerfile_resources/onomyd"])
            .await
            .stack()?;
        container_runner(&args, &[("onomyd", &dockerfile_onomyd())])
            .await
            .stack()
    }
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let daemon_home = args.daemon_home.as_ref().stack()?;
    let mut options = CosmosSetupOptions::new(daemon_home);
    options.high_staking_level = true;
    onomyd_setup(options).await.stack()?;
    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", None).await.stack()?;

    let addr = &cosmovisor_get_addr("validator").await.stack()?;
    let valoper_addr = &reprefix_bech32(addr, "onomyvaloper").stack()?;
    let tmp = sh_cosmovisor(["tendermint show-address"]).await.stack()?;
    let valcons_addr = tmp.trim();
    info!("address: {addr}");
    info!("valoper address: {valoper_addr}");
    info!("valcons address: {valcons_addr}");

    let test_deposit = token18(500.0, "anom");
    let proposal = json!({
        "title": "Text Proposal",
        "description": "a text proposal",
        "type": "Text",
        "deposit": test_deposit
    });
    cosmovisor_gov_file_proposal(daemon_home, None, &proposal.to_string(), "1anom")
        .await
        .stack()?;
    let proposals = sh_cosmovisor(["query gov proposals"]).await.stack()?;
    ensure!(proposals.contains("PROPOSAL_STATUS_PASSED"));

    // get valcons bech32 and pub key
    // cosmovisor run query tendermint-validator-set
    let valcons_set = sh_cosmovisor(["query tendermint-validator-set"])
        .await
        .stack()?;
    info!("{valcons_set}");

    // get mapping of cons pub keys and valoper addr
    // cosmovisor run query staking validators

    sh_cosmovisor_tx([format!(
        "staking delegate {valoper_addr} 1000000000000000000000anom --fees 1000000anom -y -b \
         block --from validator"
    )])
    .await
    .stack()?;
    sh_cosmovisor(["query staking validators"]).await.stack()?;

    let apr0 = get_apr_annual(valoper_addr, 6311520.0).await.stack()?;
    info!("APR: {apr0}");
    wait_for_num_blocks(1).await.stack()?;
    let apr1 = get_apr_annual(valoper_addr, 6311520.0).await.stack()?;
    info!("APR: {apr1}");
    ensure!(apr1 < apr0);
    ensure!(apr1 < 0.14);

    info!("{}", get_delegations_to(valoper_addr).await.stack()?);
    info!("{:?}", get_staking_pool().await.stack()?);
    info!("{}", get_treasury().await.stack()?);
    info!("{}", get_treasury_inflation_annual().await.stack()?);

    sh([format!(
        "cosmovisor run tx bank send {addr} onomy1a69w3hfjqere4crkgyee79x2mxq0w2pfj9tu2m 1337anom \
         --fees 1000000anom -y -b block"
    )])
    .await
    .stack()?;

    //cosmovisor run tx staking delegate onomyvaloper
    // 10000000000000000000000ibc/
    // 0EEDE4D6082034D6CD465BD65761C305AACC6FCA1246F87D6A3C1F5488D18A7B --gas auto
    // --gas-adjustment 1.3 -y -b block

    let test_crisis_denom = ONOMY_IBC_NOM;
    let test_deposit = token18(2000.0, "anom");
    cosmovisor_gov_file_proposal(
        daemon_home,
        Some("param-change"),
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
        "1anom",
    )
    .await
    .stack()?;
    wait_for_num_blocks(1).await.stack()?;
    // just running this for debug, param querying is weird because it is json
    // inside of yaml, so we will instead test the exported genesis
    sh_cosmovisor(["query params subspace crisis ConstantFee"])
        .await
        .stack()?;

    sleep(Duration::ZERO).await;
    cosmovisor_runner.terminate(TIMEOUT).await.stack()?;
    // test that exporting works
    let exported = sh_cosmovisor_no_debug(["export"]).await.stack()?;
    FileOptions::write_str("/logs/onomyd_export.json", &exported)
        .await
        .stack()?;
    let exported = yaml_str_to_json_value(&exported)?;
    ensure_eq!(
        stacked_get!(exported["app_state"]["crisis"]["constant_fee"]["denom"]),
        test_crisis_denom
    );
    ensure_eq!(
        stacked_get!(exported["app_state"]["crisis"]["constant_fee"]["amount"]),
        "1337"
    );

    Ok(())
}
