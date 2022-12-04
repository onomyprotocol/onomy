# Steps to run the full node

# How to Run an Onomy Full Node

As a Cosmos-based chain, the ONET full nodes are similar to any Cosmos full nodes. Unlike the validator flow, running a
full node requires no external software.

## Getting Started

System requirements:

- Any modern Linux distribution (RHEL 8 or Fedora 36 are preferred)
- A quad-core CPU
- 16 GiB RAM
- 320gb of storage space

## Run the full node

* Go to "deploy/scripts" folder of the repository

* Make the scripts executable

    ```
    chmod +x *
    ```

* Install chain binaries using the doc [installation](installation.md).

* Optionally expose monitoring

    ```
    ./expose-metrics.sh
    ```

* Optionally allow cors requests

    ```
    ./allow-cors.sh
    ```

* Optionally set the snapshot configuration (only if you want the node start saving snapshots)

    ```
    ./set-snapshots.sh
    ```

* Start the node

  Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

  ```
  ./start-cosmovisor-onomyd.sh &>> $HOME/.onomy/logs/onomyd.log &
  ```

  Or If you want to run the node without cosmovisor:

  ```
  ./start-onomyd.sh &>> $HOME/.onomy/logs/onomyd.log &
  ```

  Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

  ```
  ./add-service.sh cosmovisor-onomyd ${PWD}/start-cosmovisor-onomyd.sh
  ```

  Or If you want to run the node without cosmovisor:

  ```
  ./add-service.sh onomyd ${PWD}/start-onomyd.sh
  ```
