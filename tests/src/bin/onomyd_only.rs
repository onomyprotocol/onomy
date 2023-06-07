use common::container_runner;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_get_addr, cosmovisor_start, get_apr_annual, get_delegations_to,
        get_staking_pool, get_treasury, get_treasury_inflation_annual, onomyd_setup,
    },
    json_inner, onomy_std_init, reprefix_bech32,
    super_orchestrator::{
        sh,
        stacked_errors::{MapAddError, Result},
    },
    token18, yaml_str_to_json_value, Args,
};

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
        container_runner(&args, &[("onomyd", "onomyd")]).await
    }
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    onomyd_setup(daemon_home, false).await?;
    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", false, None).await?;

    let addr: &String = &cosmovisor_get_addr("validator").await?;
    let valoper_addr = &reprefix_bech32(addr, "onomyvaloper").unwrap();
    assert!((get_apr_annual(valoper_addr).await? - 13.25).abs() < 0.1);

    // make sure DAO is not delegating
    let delegations = yaml_str_to_json_value(&get_delegations_to(valoper_addr).await?)?;
    assert_eq!(
        json_inner(&delegations["delegation_responses"][0]["balance"]["amount"]),
        token18(1.0e6, "")
    );

    let staking_pool = get_staking_pool().await?;
    assert_eq!(staking_pool.bonded_tokens, 1.0e6);
    assert_eq!(staking_pool.unbonded_tokens, 0.0);

    assert!((get_treasury().await? - 100.0e6).abs() < 100.0);
    assert!((get_treasury_inflation_annual().await? - 0.13).abs() < 0.001);

    cosmovisor_runner.terminate().await?;

    Ok(())
}
