module github.com/onomyprotocol/onomy

go 1.18

require (
	github.com/cosmos/cosmos-sdk v0.45.11
	github.com/cosmos/ibc-go/v2 v2.5.0
	github.com/ethereum/go-ethereum v1.10.17
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/onomyprotocol/arc/module v0.0.0-20230307040029-595027c25e8e
	github.com/onomyprotocol/tm-load-test v0.9.1-0.20211101093435-b38e68e11c01
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	github.com/tendermint/starport v0.19.3
	github.com/tendermint/tendermint v0.34.26
	github.com/tendermint/tm-db v0.6.6
	github.com/testcontainers/testcontainers-go v0.12.0
	google.golang.org/genproto v0.0.0-20230125152338-dcaf20b6aeaa
	google.golang.org/grpc v1.52.3
	gopkg.in/yaml.v2 v2.4.0
)


replace (
	github.com/cosmos/cosmos-sdk => github.com/onomyprotocol/onomy-sdk v0.44.6-0.20230407041659-16a494c87f72
	github.com/onomyprotocol/arc/module => github.com/onomyprotocol/arc/module [FIXME]
	// or use this
	//github.com/onomyprotocol/arc/module => github.com/AaronKutch/cosmos-gravity-bridge/module v0.0.0-20230409212655-8694b0e7d8c7

	github.com/99designs/keyring => github.com/cosmos/keyring v1.2.0
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.8.1
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/jhump/protoreflect => github.com/jhump/protoreflect v1.9.0
	github.com/tendermint/tendermint => github.com/informalsystems/tendermint v0.34.26
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
