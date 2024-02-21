# How the run the bdjuno + hasura with docker compose

* Update pull updated docker images 

```
docker pull onomy/hasura-graphql-engine:latest
docker pull onomy/bdjuno:latest
docker pull onomy/big-dipper-2-mainnet:latest
```

* Start postgres

```
docker-compose up -d postgres
```

* Apply scripts from the (database/schema) repo https://github.com/onomyprotocol/bdjuno

* Create .bdjuno folder and copy config.toml (based on config-sample.toml) and genesis.json files. 
  
* Start the bdjuno

```
docker-compose up -d bdjuno
```

* Start the hasura

```
docker-compose up -d hasura
```

* Apply the hasura metadata

```
docker-compose exec hasura hasura metadata apply --admin-secret "xonomy"
```

The env HASURA_GRAPHQL_ADMIN_SECRET from the dockerfile contains the secret.

If the output is

```
INFO Metadata applied 
```

Then metadata applied. If empty or error, something went wrong.

* Start the big-deeper2 + helpers (node exporter)

```
docker-compose up -d
```