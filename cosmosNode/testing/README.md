# Scenario 1: We have two validators
 - Master and peer node have same genesis file
 - peer node has added master as persistent peer
 - both are wrking fine 
 - height is synchronising

# Scenario 2 : We have two connected validators and trying to connect third validator
  ## MASTER NODE 
 - Genesis file with information of three validators node
 - persistent peer information is already present in master node for both peer node
 - now height is fixed to 1 and not increasing
 - ERROR:attempting to connect to validator1

 ## VALIDATOR1
 - This was already connected to the Master node
 - But in it's genesis file it has information about master node and itself
 - persistent peer only has information about master node 
 - height is fixed to last value when only two validators were there
 - ERROR:attempting to connect to main validator

 ## VALIDATOR2
 - we are trying to connect this validator when two validators are already connected
 - it's genesis file has information about all three validators
 - persistent peer only has information of master node
 - now height is fixed to 1 and not increasing
 - ERROR: Error adding vote

 - Now we have tried to update the genesis file of validator 1 and restart it
 - It solves the error and all validator start runnning and their height start synchronising but 
   this process reset the last height which was when two validators were connected and start from height=1

## CONCLUSION : 
 - whenever a new validor sends a request to join the test-chain to master-node a new geneis file is created.
 - then all connected validators need to change there genesis file to new genrated genesis file and restart themselves.
 - the problem here is that the block-height is reseting to 1.
