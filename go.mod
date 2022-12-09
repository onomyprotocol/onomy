module github.com/onomyprotocol/onomy

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.45.10
	github.com/cosmos/ibc-go/v3 v3.4.0
	github.com/ethereum/go-ethereum v1.10.10
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/onomyprotocol/cosmos-gravity-bridge/module v0.0.0-20220606130009-db180e048abe
	github.com/onomyprotocol/tm-load-test v0.9.1-0.20211101093435-b38e68e11c01
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/starport v0.19.3
	github.com/tendermint/tendermint v0.34.22
	github.com/tendermint/tm-db v0.6.7
	github.com/testcontainers/testcontainers-go v0.12.0
	google.golang.org/genproto v0.0.0-20220725144611-272f38e5d71b
	google.golang.org/grpc v1.50.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
	github.com/cosmos/cosmos-sdk => github.com/onomyprotocol/onomy-sdk v0.44.6-0.20221103153534-77ffa1c3fab2
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
