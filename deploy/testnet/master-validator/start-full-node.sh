
echo "Starting onomy full node"

# Name of the onomy artifact
ONOMY=onomyd

# Todo: Check file limits before starting the node

echo "Running onomy..."
$ONOMY start --pruning=nothing

