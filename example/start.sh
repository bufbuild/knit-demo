#!/bin/bash

set -e

cd "$(dirname $0)"

export GOBIN=$(dirname $PWD)/.tmp/bin
mkdir -p $GOBIN
go install ../go/cmd/swapi-server
go install github.com/bufbuild/knit-go/cmd/knitgateway@${KNIT_GO_VERSION:=v0.1.0}

buf build buf.build/bufbuild/knit-demo \
    --type buf.knit.demo.swapi.planet.v1.PlanetService \
    --type buf.knit.demo.swapi.relations.v1.PlanetResolverService \
    -o ../.tmp/planetsvc.protoset

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

run_server "   film" $GOBIN/swapi-server -port 30481 \
    -service "buf.knit.demo.swapi.film.v1.FilmService" \
    -service "buf.knit.demo.swapi.relations.v1.FilmResolverService" &
pids="$!"
run_server " person" $GOBIN/swapi-server -port 30482 \
    -service "buf.knit.demo.swapi.person.v1.PersonService" \
    -service "buf.knit.demo.swapi.species.v1.SpeciesService" \
    -service "buf.knit.demo.swapi.relations.v1.PersonResolverService" \
    -service "buf.knit.demo.swapi.relations.v1.SpeciesResolverService" &
pids="$pids $!"
run_server " planet" $GOBIN/swapi-server -port 30483 \
    -service "buf.knit.demo.swapi.planet.v1.PlanetService" \
    -service "buf.knit.demo.swapi.relations.v1.PlanetResolverService" &
pids="$pids $!"
run_server "vehicle" $GOBIN/swapi-server -port 30484 \
    -service "buf.knit.demo.swapi.starship.v1.StarshipService" \
    -service "buf.knit.demo.swapi.vehicle.v1.VehicleService" \
    -service "buf.knit.demo.swapi.relations.v1.StarshipResolverService" \
    -service "buf.knit.demo.swapi.relations.v1.VehicleResolverService" &
pids="$pids $!"

run_server "gateway" $GOBIN/knitgateway -conf knitgateway.swapi-micro.yaml &
pids="$pids $!"

for pid in $pids; do
  wait $pid
done
pids=""
