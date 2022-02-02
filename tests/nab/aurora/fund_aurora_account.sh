#!/bin/bash
set -eux

export NEAR_ENV=local
export AURORA_ENGINE=aurora.node0
export NEAR_MASTER_ACCOUNT=aurora.node0

near call aurora.node0 fund_account '{"address": "0xBf660843528035a5A4921534E156a27e64B231fE", "amount": "23229320729235784806170624"}' \
    --accountId aurora.node0 --keyPath /root/.near-credentials/local/aurora.node0.json

echo "The aurora key is:"
cat /root/.near-credentials/local/aurora.node0.json

aurora get-balance 0xBf660843528035a5A4921534E156a27e64B231fE
aurora get-nonce 0xBf660843528035a5A4921534E156a27e64B231fE