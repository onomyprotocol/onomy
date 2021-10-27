module github.com/onomyprotocol/onomy

go 1.16

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/Microsoft/hcsshim v0.9.0 // indirect
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-20210908132153-348250c44fe7
	github.com/containerd/cgroups v1.0.2 // indirect
	github.com/cosmos/cosmos-sdk v0.42.1
	github.com/docker/docker v20.10.9+incompatible // indirect
	github.com/ethereum/go-ethereum v1.10.3
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	github.com/testcontainers/testcontainers-go v0.11.1
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
	google.golang.org/genproto v0.0.0-20211021150943-2b146023228c // indirect
	google.golang.org/grpc v1.41.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/althea-net/cosmos-gravity-bridge/module => github.com/onomyprotocol/cosmos-gravity-bridge/module v0.0.0-20211027094645-3b6b3f70239c
