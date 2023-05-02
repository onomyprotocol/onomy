module github.com/onomyprotocol/onomy

go 1.19

require (
	github.com/cosmos/cosmos-sdk v0.45.15
	github.com/cosmos/ibc-go/v2 v2.0.4
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/onomyprotocol/tm-load-test v0.9.1-0.20211101093435-b38e68e11c01
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tendermint/starport v0.19.3
	github.com/tendermint/tendermint v0.34.27
	github.com/tendermint/tm-db v0.6.6
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230125152338-dcaf20b6aeaa
	google.golang.org/grpc v1.52.3
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	// v0.45.15-ics-onomy
	github.com/cosmos/cosmos-sdk => ../cosmos-sdk
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/cometbft/cometbft v0.34.27
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
