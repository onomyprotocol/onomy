module github.com/onomyprotocol/onomy

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/cosmos/ibc-go/v2 v2.0.2
	github.com/ethereum/go-ethereum v1.10.10
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/onomyprotocol/cosmos-gravity-bridge/module v0.0.0-20220509085953-475d5774972b
	github.com/onomyprotocol/tm-load-test v0.9.1-0.20211101093435-b38e68e11c01
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/starport v0.19.3
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	github.com/testcontainers/testcontainers-go v0.12.0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e
	google.golang.org/grpc v1.45.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/onomyprotocol/onomy-sdk v0.44.6-0.20220509081037-e09065c62b6b
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
