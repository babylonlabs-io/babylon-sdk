#!/bin/bash
set -o errexit -o nounset -o pipefail
command -v shellcheck >/dev/null && shellcheck "$0"

CONTRACTS_FOLDER="cosmos-bsn-contracts"
CONTRACTS="babylon_contract btc_light_client btc_staking btc_finality"
OUTPUT_FOLDER="$(dirname "$0")/../testdata"

echo "DEV-only: copy from local built instead of downloading"

for CONTRACT in $CONTRACTS; do
  cp -f "../../${CONTRACTS_FOLDER}/artifacts/${CONTRACT}".wasm "$OUTPUT_FOLDER/"
done

cd "../../${CONTRACTS_FOLDER}"
TAG=$(git rev-parse HEAD)
cd - 2>/dev/null
echo "$TAG" >"$OUTPUT_FOLDER/version.txt"
