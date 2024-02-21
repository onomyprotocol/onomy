# State-Sync

State sync allows a new node to join a network by fetching a snapshot of the application state at a recent height
instead of fetching and replaying all historical blocks. This can reduce the time needed to sync with the network from
days to minutes.

Tendermint Core handles most of the grunt work of discovering, exchanging, and verifying state data for state sync, but
the application must take snapshots of its state at regular intervals and make these available to Tendermint via ABCI
calls, and be able to restore these when syncing a new node.

Onomy nodes can be configured to take snapshots at regular height intervals. These snapshots are stored in
$ONOMY_HOME/data/snapshots directory.

When a new node is state-synced, Tendermint will fetch a snapshot from peers in the network and provide it to the
local (empty) application, which will import it into its IAVL stores. Tendermint then verifies the application’s app
hash against the main blockchain using light client verification, and proceeds to execute blocks as usual. Note that a
state synced node will only restore the application state for the height the snapshot was taken at, and will not contain
historical data nor historical blocks.

## Setup Onomy node to take state-sync snapshots

When starting Onomy node, provide 2 more flag in the command to enable state-sync snapshots

```
onomyd start --sync.snapshot-interval 1000 --state-sync.snapshot-keep-recent 3
```

If run with the command above, onomy node will take snapshots every 1000 blocks and will keep recent 3 snapshots and
delete older ones.

## Setup a new node to sync with state-sync

When you initialize a new node with `init-full-node.sh`, script asks you if you want to set statesync. input 'y' at the
prompt and enter addresses of statesync nodes to sync up faster
