use std::time::Duration;

use common::{container_runner, dockerfile_onomyd};
use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_get_addr, cosmovisor_gov_file_proposal, cosmovisor_start, get_apr_annual,
        get_delegations_to, get_staking_pool, get_treasury, get_treasury_inflation_annual,
        sh_cosmovisor, sh_cosmovisor_no_dbg, sh_cosmovisor_tx, wait_for_num_blocks,
    },
    json_inner, onomy_std_init, reprefix_bech32,
    setups::onomyd_setup,
    super_orchestrator::{
        sh,
        stacked_errors::{MapAddError, Result},
        FileOptions,
    },
    token18, yaml_str_to_json_value, Args, ONOMY_IBC_NOM, TIMEOUT,
};
use tokio::time::sleep;

#[tokio::main]
async fn main() -> Result<()> {
    let args = onomy_std_init()?;

    if let Some(ref s) = args.entry_name {
        match s.as_str() {
            "onomyd" => onomyd_runner(&args).await,
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
        container_runner(&args, &[("onomyd", &dockerfile_onomyd())]).await
    }
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    onomyd_setup(daemon_home).await?;
    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", None).await?;

    let addr = &cosmovisor_get_addr("validator").await?;
    let valoper_addr = &reprefix_bech32(addr, "onomyvaloper").unwrap();
    info!("valoper address: {valoper_addr}");

    // make sure DAO is not delegating
    let delegations = yaml_str_to_json_value(&get_delegations_to(valoper_addr).await?)?;
    assert_eq!(
        json_inner(&delegations["delegation_responses"][0]["balance"]["amount"]),
        token18(1.0e6, "")
    );

    sh_cosmovisor_tx(
        &format!(
            "staking delegate {valoper_addr} 1000000000000000000000anom --fees 1000000anom -y -b \
             block --from validator"
        ),
        &[],
    )
    .await?;
    sh_cosmovisor("query staking validators", &[]).await?;

    info!("{}", get_apr_annual(valoper_addr).await?);

    info!("{}", get_delegations_to(valoper_addr).await?);
    info!("{:?}", get_staking_pool().await?);
    info!("{}", get_treasury().await?);
    info!("{}", get_treasury_inflation_annual().await?);
    info!("{}", get_apr_annual(valoper_addr).await?);

    wait_for_num_blocks(1).await?;
    info!("{}", get_apr_annual(valoper_addr).await?);

    sh_cosmovisor_tx(
        &format!(
            "bank send {addr} onomy1a69w3hfjqere4crkgyee79x2mxq0w2pfj9tu2m 1337anom --fees \
             1000000anom -y -b block"
        ),
        &[],
    )
    .await?;

    let test_crisis_denom = ONOMY_IBC_NOM;
    let test_deposit = token18(2000.0, "anom");
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
        "1anom",
    )
    .await?;
    wait_for_num_blocks(1).await?;
    // just running this for debug, param querying is weird because it is json
    // inside of yaml, so we will instead test the exported genesis
    sh_cosmovisor("query params subspace crisis ConstantFee", &[]).await?;

    sleep(Duration::ZERO).await;
    cosmovisor_runner.terminate(TIMEOUT).await?;
    // test that exporting works
    let exported = sh_cosmovisor_no_dbg("export", &[]).await?;
    FileOptions::write_str("/logs/onomyd_export.json", &exported).await?;
    let exported = yaml_str_to_json_value(&exported)?;
    assert_eq!(
        exported["app_state"]["crisis"]["constant_fee"]["denom"],
        test_crisis_denom
    );
    assert_eq!(
        exported["app_state"]["crisis"]["constant_fee"]["amount"],
        "1337"
    );

    Ok(())
}
