use common::container_runner;
use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_start, get_block_height, get_staking_pool, get_treasury,
        get_treasury_inflation_annual, onomyd_setup, sh_cosmovisor, wait_for_height,
    },
    nom, onomy_std_init,
    super_orchestrator::{
        sh,
        stacked_errors::{MapAddError, Result},
        STD_DELAY, STD_TRIES,
    },
    Args,
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
        container_runner(&args, &[("chain_upgrade", "onomyd")]).await
    }
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let onomy_current_version = args.onomy_current_version.as_ref().map_add_err(|| ())?;
    let onomy_upgrade_version = args.onomy_upgrade_version.as_ref().map_add_err(|| ())?;
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;
    assert_ne!(onomy_current_version, onomy_upgrade_version);
    // TODO for the next version we turn this to 'false'
    onomyd_setup(daemon_home, true).await?;
    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", false, None).await?;

    assert_eq!(
        sh_cosmovisor("version", &[]).await?.trim(),
        onomy_current_version
    );

    let upgrade_prepare_start = get_block_height().await?;
    let upgrade_height = &format!("{}", upgrade_prepare_start + 4);
    let proposal_id = "1";

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

    let description = &format!("\"upgrade {onomy_upgrade_version}\"");
    sh_cosmovisor(
        "tx gov submit-proposal software-upgrade",
        &[
            [
                onomy_upgrade_version,
                "--title",
                description,
                "--description",
                description,
                "--upgrade-height",
                upgrade_height,
            ]
            .as_slice(),
            gas_args,
        ]
        .concat(),
    )
    .await?;
    sh_cosmovisor(
        "tx gov deposit",
        &[[proposal_id, &nom(2000.0)].as_slice(), gas_args].concat(),
    )
    .await?;
    sh_cosmovisor(
        "tx gov vote",
        &[[proposal_id, "yes"].as_slice(), gas_args].concat(),
    )
    .await?;

    wait_for_height(STD_TRIES, STD_DELAY, upgrade_prepare_start + 5).await?;

    let version = sh_cosmovisor("version", &[]).await?.trim().to_owned();
    // if the build is not on a tag we get some hashed garbage on the end
    assert_eq!(version.find(onomy_upgrade_version).unwrap(), 0);

    info!("{:?}", get_staking_pool().await?);
    info!("{}", get_treasury().await?);
    info!("{}", get_treasury_inflation_annual().await?);

    cosmovisor_runner.terminate().await?;

    Ok(())
}
