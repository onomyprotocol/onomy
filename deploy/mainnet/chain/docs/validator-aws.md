# Validator on the Onomy Network

Validators should not be allowed to join set with a sufficient stake. Mainnet manager to supply NOM to those
stakeholders wishing to run the validator. This mainnet is for stakeholders only with 225k+ NOM. Anyone may run a node,
but only stakeholders may participate in validating the network.

# Steps to run the validator node

* Go to [scripts](../scripts) folder

* Make the scripts executable

    ```
    chmod +x *
    ```

* Install dependencies from source code

    ```
    ./bin.sh
    ```

* Init password store to store private keys. Run the following script to setup a password store and gpg key.

    ```
    ./init-pass.sh
    ```

Check that the output is without errors and fails.

* Set the AWS key names to use.
   ```
   export VALIDATOR_AWS_KEYS_NAME="your-full-key-name" # UPDATE THE KEY NAME 
   ```
  That key should be present on the AWS keys store

  The AWS expected keys structure for validator:
   ```
   eth_address_1_0 - the address from the eth key
   eth_private_key_1_0 - eth orchestrator private key
   mnemonic_2 - onomy eth orchestrator mnemonic
   mnemonic_1 - onomy validator mnemonic
   ``` 


* Init validator

    * Init the validator accounts from the aws keys.

    ```
    ./init-validator-keys-aws.sh
    ```

* Init full node

    ```
    ./init-full-node.sh
    ```

* Optionally with sentries

    * Start sentry nodes based on instructions from the [sentry node setup](sentry.md)

      Get the node id:
        ```
        onomyd tendermint show-node-id
        ```

      Get the node ip:

        ```
        hostname -I | awk '{print $1}'
        ```

    * Run script to set up the private connection of the validator and sentries

      Make sure to setup and start all the sentries before running this script as you will need to provide IPs of all
      the sentry nodes.

        ```
        ./set-sentry.sh
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

* Then check the balance on validator node

    ```
    onomyd q bank balances $(onomyd keys show validator -a --keyring-backend pass)
    ```

  If the "amount" of noms >= 250k is updated you are ready to become a validator

* Create a new onomy validator

  **!!! Attention !!!**

  Once you create a validator, you will have 500 blocks (approximately 40 minutes) to start the orchestrator, otherwise
  the validator will be slashed. Read the full instruction before proceed with the next step.

    ```
    ./create-validator.sh
    ```

Also you can check all current validators now.

    ```
    onomyd q staking validators
    ```

* Optionally if it isn't funded send 1k noms tokens from you validator to your orchestrator

    ```
    onomyd tx bank send $(onomyd keys show validator -a --keyring-backend pass) $(onomyd keys show eth-orchestrator -a --keyring-backend pass) 1000000000000000000000anom --gas auto --gas-adjustment 1.5 --chain-id=onomy-mainnet-1 --keyring-backend pass
    ```

Check the orchestrator balance

    ```
    onomyd q bank balances $(onomyd keys show eth-orchestrator -a --keyring-backend pass)
    ```

* Run orchestrator

    * Init eth orchestrator

      ```
      ./init-eth-orchestrator-aws.sh
      ```
    * Check that your Ethereum address is in the list of current valset

      ```
      onomyd q gravity current-valset
      ```

    * Before running the script set env variable

        * ETH_RPC_ADDRESS - the RPC address of the Ethereum node
        
    * start-orchestrator
    
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
