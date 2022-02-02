module github.com/onomyprotocol/onomy

go 1.16

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/Microsoft/hcsshim v0.9.0 // indirect
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-20210908132153-348250c44fe7
	github.com/containerd/cgroups v1.0.2 // indirect
	github.com/cosmos/cosmos-sdk v0.42.6
	github.com/docker/docker v20.10.9+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/onomyprotocol/near-aurora-bridge/module v0.0.0-20220130125415-880c96c08165
	github.com/onomyprotocol/tm-load-test v0.9.1-0.20211101093435-b38e68e11c01
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/spm v0.1.2
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	github.com/testcontainers/testcontainers-go v0.11.1
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5 // indirect
	google.golang.org/grpc v1.43.0 // indirect
	google.golang.org/protobuf v1.27.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/althea-net/cosmos-gravity-bridge/module => github.com/onomyprotocol/cosmos-gravity-bridge/module v0.0.0-20211125165615-060eb3403d3c

replace github.com/cosmos/cosmos-sdk => github.com/onomyprotocol/onomy-sdk v0.42.10-0.20211228140704-1a3046991600
