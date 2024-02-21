# Steps to run the ethereum node

* Go to [scripts](../scripts) folder

* Make the scripts executable

    ```
    chmod +x *
    ```

* Install dependencies from source code

    ```
    ./bin.sh
    ```

* Start the node

  Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

    ```
    ./start-geth.sh &>> $HOME/.onomy/logs/geth.log &
    ```

  Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

    ```
    ./add-service.sh geth ${PWD}/start-geth.sh
    ```

* Optionally run node exporter

    ```
    ./start-node-exporter.sh &>> $HOME/.onomy/logs/node-exporter.log &
    ```

  Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

    ```
    ./add-service.sh node-exporter ${PWD}/start-node-exporter.sh
    ```