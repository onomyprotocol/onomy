
echo "Starting onomy full node and faucet"

# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9090"
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"
# Onomy chain demons
STAKE_DENOM="nom"

echo "Running onomy..."
$ONOMY start --pruning=nothing &

echo "Waiting $ONOMY_NODE_NAME to launch gRPC $ONOMY_GRPC_PORT..."

while ! timeout 1 bash -c "</dev/tcp/$ONOMY_HOST/$ONOMY_GRPC_PORT"; do
  sleep 1
done

echo "$ONOMY_NODE_NAME launched"

# ------------------ Run faucet ------------------
ONOMY_ORCHESTRATOR_NAME=$(jq -r .name $ONOMY_HOME/orchestrator_key.json)

echo "Starting faucet based on validator account"
faucet -cli-name=$ONOMY -account-name="$ONOMY_ORCHESTRATOR_NAME" $ONOMY_KEYRING_FLAG --denoms $STAKE_DENOM credit-amount=100000000  -max-credit=200000000 &
