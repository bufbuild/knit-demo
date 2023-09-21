# SWAPI Microservices Example

This directory contains example configuration and scripts for running the SWAPI demo app
as if it were a set of microservices. This demonstrates how the Knit Gateway can be
configured to operate in such an environment.

* Run the `start.sh` script to build and start all of the server processes. This script
  will build the `swapi-server` program in this repo and also install `knitgateway`.
  It then starts five different server processes:
  1. _film_: This server provides the API for films. It also provides resolver RPCs
     for resolving references to films. It is an instance of `swapi-server`, running
     on port 30481.
  2. _person_: This server provides the API for people and species, as well as the
     resolver RPCs for person and species references. It is also an instance of
     `swapi-server`, and it runs on port 30482.
  3. _planet_: This server provides the APIs for planets. This instance of `swapi-server`
     runs on port 30483.
  4. _vehicle_: This server provides the APIs for vehicles and starships. This is
     the final instance of `swapi-server`, running on port 30484.
  5. _gateway_: This is the `knitgateway`, which processes Knit queries and handles
     dispatching RPCs to the above four servers.

* The gateway service started by the above script uses the `knitgateway.swapi-micro.yaml`
  configuration file. This is what wires up the gateway to communicate with the other
  four "microservices".

* The example client script in the `ts` folder of this repo works when the above
  servers are running. When you execute the example query, you can see the output
  of all of the servers as the query components are dispatched by the gateway.
