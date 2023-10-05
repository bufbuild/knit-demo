#!/bin/bash

set -e

cd "$(dirname $0)"

export GOBIN=$PWD/.tmp/bin
mkdir -p $GOBIN
go install ./go/cmd/swapi-server
go install github.com/bufbuild/knit-go/cmd/knitgateway@v0.1.0

curl -o ./.tmp/knitgateway.yaml https://raw.githubusercontent.com/bufbuild/knit-go/main/knitgateway.example.yaml

function cleanup() {
  for pid in $pids; do
    kill $pid 2>/dev/null || true
  done
}

trap cleanup EXIT

function run_server() {
  server_name="$1"
  shift
  exec > >(trap "" INT TERM; sed 's/^/'"$server_name"': /')
  exec 2> >(trap "" INT TERM; sed 's/^/'"$server_name"': /' >&2)
  exec "$@"
}

run_server "  swapi" $GOBIN/swapi-server &
pids="$!"

run_server "gateway" $GOBIN/knitgateway -conf ./.tmp/knitgateway.yaml &
pids="$pids $!"

# We want to make sure above servers are up and running before we
# run the next step. So give it a second (literally).
sleep 1

cd ts
npm install
npm run start
