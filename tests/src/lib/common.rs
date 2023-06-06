use onomy_test_lib::{
    super_orchestrator::{
        docker::{Container, ContainerNetwork},
        sh,
        stacked_errors::Result,
    },
    Args, TIMEOUT,
};

/// Useful for running simple container networks that have a standard format and
/// don't need extra build or volume arguments.
pub async fn container_runner(
    args: &Args,
    dockerfiles_and_entry_names: &[(&str, &str)],
) -> Result<()> {
    let bin_entrypoint = &args.bin_name;
    let container_target = "x86_64-unknown-linux-gnu";
    let logs_dir = "./tests/logs";

    // build internal runner
    sh("cargo build --release --bin", &[
        bin_entrypoint,
        "--target",
        container_target,
    ])
    .await?;

    let mut cn = ContainerNetwork::new(
        "test",
        dockerfiles_and_entry_names
            .iter()
            .map(|(dockerfile, entry_name)| {
                Container::new(
                    entry_name,
                    Some(&format!("./tests/dockerfiles/{dockerfile}.dockerfile")),
                    None,
                    &[(logs_dir, "/logs")],
                    Some(&format!(
                        "./target/{container_target}/release/{bin_entrypoint}"
                    )),
                    &["--entry-name", entry_name],
                )
            })
            .collect(),
        // TODO
        true,
        logs_dir,
    )?;
    cn.run_all(true).await?;
    cn.wait_with_timeout_all(true, TIMEOUT).await.unwrap();
    Ok(())
}
