#!/bin/bash
set -o nounset -o pipefail
command -v shellcheck >/dev/null && shellcheck "$0"

OWNER="babylonlabs-io"
REPO="babylon-contract"
CONTRACTS="babylon_contract btc_staking"
OUTPUT_FOLDER="$(dirname "$0")/../testdata"

[ $# -ne 1 ] && echo "Usage: $0 <version>" && exit 1
type curl >&2

TAG="$1"

for CONTRACT in $CONTRACTS
do
  echo -n "Downloading $CONTRACT..." >&2
  URL="https://github.com/$OWNER/$REPO/releases/download/$TAG/$CONTRACT.wasm.zip"
  curl -s -L -H 'Accept: application/octet-stream' "$URL" >"$OUTPUT_FOLDER/$CONTRACT.wasm.zip"
  unzip -o "$OUTPUT_FOLDER/$CONTRACT.wasm.zip"
  rm -f "$OUTPUT_FOLDER/$CONTRACT.wasm.zip"
  echo "done." >&2
done
echo "$TAG" >"$OUTPUT_FOLDER/version.txt"
