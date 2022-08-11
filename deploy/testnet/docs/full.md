# Steps to run the full node

* Go to [scripts](../scripts) folder

* Make the scripts executable

    ```
    chmod +x *
    ```

* Install dependencies from source code

    ```
    ./bin.sh
    ```

* Init full node

    ```
    ./init-full-node.sh
    ```

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

  Before running the script please set up "ulimit > 65535"

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
