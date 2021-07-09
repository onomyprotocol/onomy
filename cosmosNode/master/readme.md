## Steps to build cosmos node
1. create docker image using this docker file
2. using this image run container and command for create container will be
 
if you want to run docker container using docker service then command will be
   
docker service create --replicas 1 -t -d --name=[service name] -p 26656:26656 -p 26657:26657 -p 1317:1317 -p 9090:9090 [docker image name] gravity --home /root/.gravity/ --address tcp://0.0.0.0:26655 --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090 --log_level error --p2p.laddr tcp://0.0.0.0:26656 --rpc.pprof_laddr 0.0.0.0:6060 start

and if you want to run using docker run command then 

docker run -d -p 26658:26658 -p 26659:26659 -p 1318:1318 -p 9091:9091 [docker image name] gravity --home /root/.gravity/ --address tcp://0.0.0.0:26655 --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090 --log_level error --p2p.laddr tcp://0.0.0.0:26656 --rpc.pprof_laddr 0.0.0.0:6060 start



