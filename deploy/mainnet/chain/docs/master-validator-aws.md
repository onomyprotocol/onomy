~~# Steps to run the master node

* Go to [scripts](../scripts) folder

* Make the scripts executable

    ```
    chmod +x *
    ```

* Install dependencies from source code

    ```
    ./bin-master.sh
    ```

* Init password store to store private keys. Run the following script to setup a password store
  and gpg key.

    ```
    ./init-pass.sh
    ```
  Check that the output is without errors and fails.

* Set the AWS key names to use.
   ```
   export VALIDATOR_AWS_KEYS_NAME="validator_1_keys"
   export CONTRACT_DEPLOYER_AWS_KEYS_NAME="ETH-Deployer"
   ```
  That key should be present on the AWS keys store

  The AWS expected keys structure for validator:
   ```
   eth_address_1_0 - the address from the eth key
   eth_private_key_1_0 - eth orchestrator private key
   mnemonic_2 - onomy eth orchestrator mnemonic
   mnemonic_1 - onomy validator mnemonic
   ``` 

  The AWS expected keys structure for contract deployer
  ```
  eth_private_key_0 - the eth private key to deploy contracts
  ```

* Add genesis accounts

  This is the accounts doc with the generated CLI inputs:
  [Accounts spreadsheets.](https://docs.google.com/spreadsheets/d/1zxX5Wx6PoX21r5GWprCzSLMJo0CXF1gWFdKXlJ2BKiw/edit#gid=354042563)

    * Take the "Account lines" from the "Basic Accounts" tab and save to "../genesis/accounts.txt"
      Check that tab contains all validators accounts as well.

    * Take the "Vesting lines" from the "Vesting Accounts" tab and save to "../genesis/accounts-vesting.txt"

* Set orchestrator ETH address
    ```
    export ONOMY_ETH_ORCHESTRATOR_VALIDATOR_ADDRESS=$(aws secretsmanager get-secret-value --secret-id $VALIDATOR_AWS_KEYS_NAME | jq --raw-output '.SecretString' | jq -r '.eth_address_1_0')
    ```
  This is the Ethereum public address with which the orchestrator will be running.

* Run init script

  In case you use AWS
    ```
    ./init-master-aws.sh
    ```

* Start the node to deploy required contracts

  Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

    ```
    ./start-onomyd.sh &>> $HOME/.onomy/logs/onomyd.log &
    ```
  Check that logs are without errors.

* Deploy gravity contract

    * Before running the script set env variable:

        * ETH_RPC_ADDRESS - the RPC address of the Ethereum node

    * Deploy eth bridge contract

      ```
      ./deploy-and-set-eth-bridge-contract-aws.sh
      ```

    * Copy the "assets/bridge/addresses.json" file to the "environments" repo.
      !!! This step is very important, since it provides the bridge addresses for the validators !!!

* Run the orchestrator to deploy the anom representation
  It is better to set different orchestrator and deployer keys.

    * Init eth orchestrator
    ```
    ./init-master-eth-orchestrator.sh
    ```

    * Before running the script set env variables (or add to the script):

        * ETH_RPC_ADDRESS - the RPC address of the Ethereum node

    * Start master orchestrator
      ```
      ./start-eth-orchestrator-aws.sh &>> $HOME/.onomy/logs/eth-orchestrator.log &
      ```

      Check that logs are without errors. And await for the message:
      ```
      Oracle resync complete, Oracle now operational 
      ```

* Deploy the anom representation and set the mapping the genesis params (it might take more than 10 mins, because of the
  orchestrator block delay)

  ```
  ./deploy-anom-representation-and-set-valset-rewards-aws.sh
  ```

* Stop the orchestrator
  ```
  ./stop-eth-orchestrator.sh
  ```

* Stop onomy node
  ```
  ./stop-onomyd.sh
  ```

* Check the logs to be sure that both bridge and onomy are stopped

* Copy genesis to environments repo

  Copy the "$HOME/.onomy/config/genesis.json" file to the "../genesis/genesis-mainnet-1.json"
  in the environments repo to be used by other scripts.

* Purge the temporary generated data and logs
   ```
   ./purge-onomy-data-and-logs.sh
   ```

* Get the node id:

    ```
    onomyd tendermint show-node-id
    ```

* Get the node ip:

    ```
    hostname -I | awk '{print $1}'
    ```

* Start initial seed node(s) based on instructions from the [seed](seed.md)

* Start sentry node(s) based on instructions from the [sentry](sentry.md)

* Run script to set up the private connection of the validator and sentries

    ```
    ./set-sentry.sh
    ```

* Expose onomy node monitoring

    ```
    ./expose-metrics.sh
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

* Run orchestrator

    * Before running the script set env variables (or add to the script):

        * ETH_RPC_ADDRESS - the RPC address of the Ethereum node

    * Start master orchestrator (the keys will be resuld from the prev steps)

      ```
      ./start-eth-orchestrator-aws.sh &>> $HOME/.onomy/logs/eth-orchestrator.log &
      ```      

      Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

      ```
      ./add-service.sh eth-orchestrator ${PWD}/start-eth-orchestrator-aws.sh
      ```

      To get the orchestrator logs you can use the command:
      ```
        journalctl -u eth-orchestrator.service -n 100 --no-pager
      ```
      If in the last lines you see the message ``` Orchestrator resync complete, Oracle now operational``` then the
      orchestrator successfully restarted.


* Optionally run node exporter

    ```
    ./start-node-exporter.sh &>> $HOME/.onomy/logs/node-exporter.log &
    ```

  Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.

    ```
    ./add-service.sh node-exporter ${PWD}/start-node-exporter.sh
    ```

