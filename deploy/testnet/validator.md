# Steps to run the peer-validator node

* Install dependencies from source code

```
./bin.sh
```

* Init full node

```
./init-full-node.sh
```

* Optionally with sentries

    * Start sentry nodes based on instructions from the [sentry](../sentry/readme.md)

        Get the node id:
        ```
        onomyd tendermint show-node-id
        ```
        
        Get the node ip:
        
        ```
        hostname -I | awk '{print $1}'
        ```
        
    * Run script to set up the private connection of the validator and sentries

    You will need to provide the sentry IPs.
    
    ```
    ./set-sentry.sh
    ```

* Expose onomy node monitoring

```
./expose-metrics.sh
```

* Start the node

Before run the script please set up "ulimit > 65535" ([example-hrel8](set-ulimit-hrel8.md))

```
./start-node.sh
```

* Init validator

Init the validator account to be deposited

```
./init-validator.sh
```

* Get tokens from master account
    * Request tokens (for validator) from the master account by "text" request to onomy.
  
    * Then check the balance on validator node

    ```
    onomyd q bank balances {validator-address}
    ```

    If the "amount" is updated you are ready to become a validator

* Create a new onomy validator

```
./create-validator.sh
```

Also you can check all current validators now.

```
onomyd q staking validators
```

* Send some tokens from you validator to your orchestrator

```
onomyd tx bank send {validator-address} {orchestrator-address} 5000000000000000000anom --chain-id=onomy-testnet --keyring-backend test
```

Check the orchestrator balance now

```
onomyd q bank balances {orchestrator-address}
```

* Init gbt

```
./init-gbt.sh
```

Check that your Ethereum address is in the list of curren valset

```
onomyd q gravity current-valset
```

* Run orchestrator

Before run the script please set env variable:

* ETH_RPC_ADDRESS - the RPC address of the Ethereum node

```
./start-orchestrator.sh
```

* Run node exporter

```
./start-node-exporter.sh
```

## Setup auto-start

Add to your crontab or /etc/init.d scripts:

* `start-node.sh`
* `start-orchestrator.sh` (set envs or updated script with ETH_RPC_ADDRESS envs)
* `start-node-exporter.sh`

***If you used the bin.sh installation and what to use the scripts for the auto-start, additionally you need to
add ```export PATH=$PATH:$ONOMY_HOME/bin``` to your scripts after the ```ONOMY_HOME=$HOME/.onomy```***
