use onomy_test_lib::{
    dockerfiles::onomy_std_cosmos_daemon,
    super_orchestrator::{
        docker::{Container, ContainerNetwork, Dockerfile},
        sh,
        stacked_errors::{Result, StackableErr},
    },
    Args, TIMEOUT,
};

pub const ONOMYD_VERSION: &str = "v1.1.4";

pub fn dockerfile_onomyd() -> String {
    onomy_std_cosmos_daemon("onomyd", ".onomy", ONOMYD_VERSION, "onomyd")
}

/// Useful for running simple container networks that have a standard format and
/// don't need extra build or volume arguments.
pub async fn container_runner(args: &Args, name_and_contents: &[(&str, &str)]) -> Result<()> {
    let logs_dir = "./tests/logs";
    let dockerfiles_dir = "./tests/dockerfiles";
    let bin_entrypoint = &args.bin_name;
    let container_target = "x86_64-unknown-linux-gnu";

    // build internal runner
    sh([
        "cargo build --release --bin",
        bin_entrypoint,
        "--target",
        container_target,
    ])
    .await?;

    let mut containers = vec![];
    for (name, contents) in name_and_contents {
        containers.push(
            Container::new(name, Dockerfile::contents(contents))
                .external_entrypoint(
                    format!("./target/{container_target}/release/{bin_entrypoint}"),
                    ["--entry-name", name],
                )
                .await
                .stack()?,
        );
    }

    let mut cn =
        ContainerNetwork::new("test", containers, Some(dockerfiles_dir), true, logs_dir).stack()?;
    cn.add_common_volumes([(logs_dir, "/logs")]);
    let uuid = cn.uuid_as_string();
    cn.add_common_entrypoint_args(["--uuid", &uuid]);
    cn.run_all(true).await.stack()?;
    cn.wait_with_timeout_all(true, TIMEOUT).await.stack()?;
    cn.terminate_all().await;
    Ok(())
}
