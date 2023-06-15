use clap::Parser;
use onomy_test_lib::super_orchestrator::{
    ctrlc_init, docker_helpers::auto_exec_i, stacked_errors::Result, std_init,
};

/// Runs auto_exec_i
#[derive(Parser, Debug)]
#[command(about)]
struct Args {
    /// Name of the container
    #[arg(short, long)]
    container_name: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    std_init()?;
    ctrlc_init()?;
    let args = Args::parse();
    auto_exec_i(&args.container_name).await?;
    Ok(())
}
