# Steps to run the sentry node

* Go to [scripts](../scripts) folder

* Make the scripts executable

    ```
    chmod +x *
    ```

* Install dependencies from source code

    ```
    ./bin.sh
    ```

* Init sentry node

    ```
    ./init-sentry-node.sh
    ```

* Init statesync or use [genesis binaries](genesis-binaries.md) instruction to run from the genesis block.

    ```
    ./init-statesync.sh
    ```

* Optionally expose monitoring

    ```
    ./expose-metrics.sh
    ```

* Start the node

  Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

    ```
    ./start-cosmovisor-onomyd.sh &>> $HOME/.onomy/logs/onomyd.log &
    ```

  Or If you want to run the node without cosmovisor (not supported by the genesis binaries):

    ```
    ./start-onomyd.sh &>> $HOME/.onomy/logs/onomyd.log &
    ```

  Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

    ```
    ./add-service.sh cosmovisor-onomyd ${PWD}/start-cosmovisor-onomyd.sh
    ```

  Or If you want to run the node without cosmovisor (not supported by the genesis binaries):

    ```
    ./add-service.sh onomyd ${PWD}/start-onomyd.sh
    ```

* Get the node ip to share with the validator node:

    ```
    hostname -I | awk '{print $1}'
    ```

* Optionally run node exporter

    ```
    ./start-node-exporter.sh &>> $HOME/.onomy/logs/node-exporter.log &
    ```

Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

    ```
    ./add-service.sh node-exporter ${PWD}/start-node-exporter.sh
    ```
