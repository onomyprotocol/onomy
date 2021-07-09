# Full Node

Full node is a node which just joins the ntework but do not participate in consensus. Also it is not a validator node, it just syncs with all other nodes present in the network

Here we'll see how you can create a full node and connect it with the gravity-testnet.

## Docker File

This docker file is used to create  docker image which when run can start a peer validator node.


## start.sh

Currently we are running this file after going inside our container, it clones the config branch which holds a folder having genesis.json file and seed file which holds the peer info than it change that genesis file with full-node genesis file and add seed info to full node. 



