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

* Init validator

    * Init the validator account to be deposited. The script will just show the keys if they already exist.

    ```
    ./init-validator-keys.sh
    ```

    * Or manually get the existing validator address

    ```
    onomyd keys show validator -a --keyring-backend pass
    ```

  **!!! Attention !!!***

  In case you want the address to be added to the genesis file (before the first launch) share the validator address
  with the onomy team. And wait until the genesis file is ready. Then you will need to update/pull the "chain" folder
  and proceed with the next steps.

* Check that pass is loaded

When the genesis is ready, and your node is ready as well, check that the required key is in the pass, run
```
onomyd keys show validator -a --keyring-backend pass
```
If you see the validator address, then proceed with the "Init full node".

But if you see an error, like this:
```
onomyd keys show validator -a --keyring-backend pass
gpg: decryption failed: No secret key
Error: validator is not a valid name or address: exit status 2
```
Then run the
```
pass keyring-onomy/validator.info
```
And then repeat the previously failed operation.

* Init full node

    ```
    ./init-full-node.sh
    ```

* Init statesync or use [genesis binaries](genesis-binaries.md) instruction to run from the genesis block.

    ```
    ./init-statesync.sh
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

* Get tokens from master account (***!!! Only in case your account wasn't set to the genesis file. If it was set - skip
  that step !!!***).
    * Request 250k noms (225k min stake for validator, 1k for orchestrator and 1k for transactions payments) tokens for
      validator from the master account by "text" request to onomy. Example how to send
    ```
    onomyd tx bank send {sender-address} {validator-address} 250000000000000000000000anom --gas auto --gas-adjustment 1.5 --chain-id=onomy-mainnet-1 --keyring-backend pass
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

* Send 1k noms tokens from you validator to your orchestrator

    ```
    onomyd tx bank send $(onomyd keys show validator -a --keyring-backend pass) $(onomyd keys show eth-orchestrator -a --keyring-backend pass) 1000000000000000000000anom --gas auto --gas-adjustment 1.5 --chain-id=onomy-mainnet-1 --keyring-backend pass
    ```

Check the orchestrator balance now

    ```
    onomyd q bank balances $(onomyd keys show eth-orchestrator -a --keyring-backend pass)
    ```

* Run orchestrator

    * Also add the private key of the ethereum wallet which will be used for orchestrator

      ```
      pass insert keyring-onomy/eth-orchestrator-eth-private-key  
      ```

    * Init eth orchestrator

      ```
      ./init-eth-orchestrator.sh
      ```
    * Check that your Ethereum address is in the list of current valset

      ```
      onomyd q gravity current-valset
      ```

    * Before running the script set env variable

        * ETH_RPC_ADDRESS - the RPC address of the Ethereum node

    * start-orchestrator

       ```
       ./start-eth-orchestrator.sh &>> $HOME/.onomy/logs/eth-orchestrator.log &
       ```

      Or add and start as a service (strongly recommended). You need to run it from the **sudo** user.
       ```
        ./add-service.sh eth-orchestrator ${PWD}/start-eth-orchestrator.sh
       ```

      !!! If you use the service mode, then after the reboot the service will request the orchestrator key from the
      pass, which is protected by the "Passphrase" and will fail. !!!

      In order to restore the orchestrator after the reboot (or in the case of `gpg: decryption failed: No secret key` error), you need to call
       ```
        ./load-keys.sh
       ```
      The script will load the required key to the temp cache and the orchestrator service will be able to get it after
      the next automatic restart.

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


# What is Jailing

When a validator disconnects from the network due to connection loss or server fail or it double signs, it needs to be
eliminated from the validator list. It is known as 'jailing'. A validator is jailed if it fails to validate at least 50%
of the last 100 blocks.

When jailed due to downtime, the validator's total stake is slashed by 1%. and if jailed due to double signing,
validaor's total stake is slashed by 5%.

Once jailed, validators can be unjailed again after 10 minutes. These configurations can be found in the genesis file
under the slashing section

```
"slashing": {
      "params": {
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600s",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": [],
      "missed_blocks": []
    }
```

### Unjailing validator

In order to unjail the validator, you may run the following command once 10 minutes have passed

```
onomyd tx slashing unjail --from <validator-name> --chain-id=onomy-mainnet-1 --gas auto --gas-adjustment 1.5 --keyring-backend pass
```


