use log::info;
use onomy_test_lib::{
    cosmovisor::{
        cosmovisor_gov_proposal, cosmovisor_start, get_block_height, get_staking_pool,
        get_treasury, get_treasury_inflation_annual, sh_cosmovisor, wait_for_height,
    },
    nom, onomy_std_init,
    setups::onomyd_setup,
    super_orchestrator::{
        docker::{Container, ContainerNetwork, Dockerfile},
        sh,
        stacked_errors::{MapAddError, Result},
        STD_DELAY, STD_TRIES,
    },
    Args, TIMEOUT,
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
        container_runner(&args).await
    }
}

async fn container_runner(args: &Args) -> Result<()> {
    let logs_dir = "./tests/logs";
    let dockerfiles_dir = "./tests/dockerfiles";
    let bin_entrypoint = &args.bin_name;
    let container_target = "x86_64-unknown-linux-gnu";

    // build internal runner
    sh("cargo build --release --bin", &[
        bin_entrypoint,
        "--target",
        container_target,
    ])
    .await?;

    let mut cn = ContainerNetwork::new(
        "test",
        vec![Container::new(
            "onomyd",
            Dockerfile::Path(format!("{dockerfiles_dir}/chain_upgrade.dockerfile")),
            Some(&format!(
                "./target/{container_target}/release/{bin_entrypoint}"
            )),
            &["--entry-name", "onomyd"],
        )],
        None,
        true,
        logs_dir,
    )?
    .add_common_volumes(&[(logs_dir, "/logs")]);
    cn.run_all(true).await?;
    cn.wait_with_timeout_all(true, TIMEOUT).await.unwrap();
    Ok(())
}

async fn onomyd_runner(args: &Args) -> Result<()> {
    let onomy_current_version = args.onomy_current_version.as_ref().map_add_err(|| ())?;
    let onomy_upgrade_version = args.onomy_upgrade_version.as_ref().map_add_err(|| ())?;
    let daemon_home = args.daemon_home.as_ref().map_add_err(|| ())?;

    info!("current version: {onomy_current_version}, upgrade version: {onomy_upgrade_version}");

    onomyd_setup(daemon_home).await?;
    let mut cosmovisor_runner = cosmovisor_start("onomyd_runner.log", None).await?;

    assert_eq!(
        sh_cosmovisor("version", &[]).await?.trim(),
        onomy_current_version
    );

    let upgrade_prepare_start = get_block_height().await?;
    let upgrade_height = &format!("{}", upgrade_prepare_start + 4);

    let description = &format!("\"upgrade {onomy_upgrade_version}\"");

    cosmovisor_gov_proposal(
        "software-upgrade",
        &[
            onomy_upgrade_version,
            "--title",
            description,
            "--description",
            description,
            "--upgrade-height",
            upgrade_height,
        ],
        &nom(2000.0),
        "1anom",
    )
    .await?;

    wait_for_height(STD_TRIES, STD_DELAY, upgrade_prepare_start + 5).await?;

    assert_eq!(
        sh_cosmovisor("version", &[]).await?.trim(),
        onomy_upgrade_version
    );

    info!("{:?}", get_staking_pool().await?);
    info!("{}", get_treasury().await?);
    info!("{}", get_treasury_inflation_annual().await?);

    cosmovisor_runner.terminate(TIMEOUT).await?;
    Ok(())
}
