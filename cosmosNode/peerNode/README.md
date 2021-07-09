# Peer-Validator Node

Here we'll se how you can create a validator node and connect it with the gravity-testnet.

## Docker File

This docker file is used to create  docker image which when run can start a peer validator node.


## start.sh

It will run with docker file when gitAction starts and it will copy the gentx folder of the machine image to a gentx folder and push that folder to GitHub. 


## update-peer-node.sh

Currently we are running this file after going inside our container, it clones the config branch which holds a folder having genesis.json file and seed file which holds the peer info than it change that genesis file with validator genesis file and add a persistent peer info to validator. 
