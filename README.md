knit-demo
=========

`knit-demo` is an example RPC service and gateway configuration built with
[Knit](https://github.com/bufbuild/knit).

Its API is defined by a [Protocol Buffer schema](https://buf.build/bufbuild/knit-demo),
and the service implementation serves the same data as is available using
[The Star Wars API](https://swapi.dev).

## Using the Demo Server

This example application, both the RPC services and the Knit gateway, is
available to interact with at https://knit-demo.connect.build/.

You can use the sources in the [`ts`](https://github.com/bufbuild/knit-demo/tree/main/ts)
folder of this repo as a client. You can edit the `index.ts` file to experiment with
different queries.

### Generated Code

The generated TypeScript code, both for using Connect and Knit, are available in this repo at
[https://github.com/bufbuild/knit-demo/tree/main/ts/gen](https://github.com/bufbuild/knit-demo/tree/main/ts/gen)

The generated Go code can be found in this repo at
[https://github.com/bufbuild/knit-demo/tree/main/go/gen](https://github.com/bufbuild/knit-demo/tree/main/go/gen)

### Building and Running

If you want to build the demo from source, you can do so using the `go` tool
and the sources in this repo:
```
go install github.com/bufbuild/knit-demo/go/cmd/swapi-server
```

The `go/cmd/swapi-server` program is the backend Connect server for the API.
It provides the same data as the API at https://swapi.dev. Once built/installed,
you can run this without any command-line flags. By default, it listens on port 30485.

To access the Star Wars API via a Knit client, there are two options:
* You can then run the `knitgateway` in the [`knit-go`](https://github.com/bufbuild/knit-go/tree/main/cmd/knitgateway)
  repo using the `knitgateway.example.yaml` config file in that repo. You will then have
  a Knit server on port 30480 that provides the Star Wars API.
* You can instead run `swapi-server` with the `--embed-gateway` flag. With this flag, both
  the Connect API and the Knit Gateway run in the same process and on the same port. So
  you will have a Knit server on port 30485.

## Status: Alpha

Knit is undergoing initial development and is not yet stable.

## Legal

Offered under the [Apache 2 license](https://github.com/bufbuild/knit-demo/blob/main/LICENSE).

