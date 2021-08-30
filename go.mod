module github.com/onomyprotocol/onomy

go 1.16

require (
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-20210623144132-d71cc5bf08f4
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cosmos/cosmos-sdk v0.42.1
	github.com/ethereum/go-ethereum v1.10.3
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/althea-net/cosmos-gravity-bridge/module => github.com/onomyprotocol/cosmos-gravity-bridge/module v0.0.0-20210716023927-0fb1bd34390a
