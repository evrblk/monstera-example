# Monstera Examples

[![Go](https://github.com/evrblk/monstera-example/actions/workflows/go.yml/badge.svg)](https://github.com/evrblk/monstera-example/actions/workflows/go.yml)

Several examples of how to build applications with [Monstera framework](https://github.com/evrblk/monstera). They are
based on popular System Design type interview questions:

* `dlocks` - Distributed RW-locks
* `ledger` - Bookkeeping service for account balances and transactions
* `tinyurl` - URL shortening service

All those projects follow the same structure which is described in details below. Some project-specific details
can be found in corresponding `README.md` files.

__Here is a bare minimum of docs you must read before jumping into this codebase:__

* https://everblack.dev/docs/monstera/overview/
* https://everblack.dev/docs/monstera/units-of-work/

Monstera framework does not force any particular application core implementation, method routing mechanism, or any 
specific wire format. It is up to you to define that. However, over time I developed a certain style of how all
Everblack services are implemented. To separate a clean part of the framework from that opinionated part I made
two packages: `github.com/evrblk/monstera` for the core part, and `github.com/evrblk/monstera/x` for the rest.
However, a lot of things are not generalizable or extractable into a library. And those example applications show 
how all of them can be assembled together.

## Application cores

All cores are implemented in my opinionated way and serve as an example of how it can be done. You are free to
do it any way you want, with different in-memory data structures or other embedded databases.

In most of those examples, Application cores store data in BadgerDB. There is one instance of BadgerDB per process, 
so multiple shards and multiple cores share it. To avoid conflicts, each table is prefixed with table IDs (in 
`tables.go`). Each shard has its own boundaries (`lowerBound` and `upperBound`). Take a look how keys are built 
for tables and indexes (typically in the bottom of files with application cores and inside `monstera/x` package too).

All core data structures are defined in protobufs in `corepb/*`. Those structures are exposed from Monstera stubs
and used by application cores to store data in BadgerDB. `corepb/cloud.proto` has high level containers for requests
and responses that are actually passed by Monstera. Monstera does not know anything about implementation of your
application cores and only passes binary blobs as requests and responses for reads and updates. Message routing to a
binary blob and from a blob is based on `oneof` protobuf structure (see `adapters.go` and `stubs.go`).

Take a look at unit tests. Application cores are easily testable without any mocks, and even very complex business 
logic can be tested by feeding the correct seqeunce of commands since all application cores are state machines without 
side effects.

## Gateway server

A gateway (or frontend) server is the public API part of the system. In this example, it serves gRPC, but it can be
anything you want (OpenAPI, ConnectRPC, gRPC, Gin, etc). Gateway gRPC is defined in `gatewaypb/*`. Protos are not
shared between gateway and core parts for clean separation of core business layer and presentation layer. The code for
converting between them lives in `pbconv.go`. `server.go` is the implementation of the gateway API. It is the entry 
point for all user actions, and if you want to trace and understand the lifecycle of a request then start from here.

Gateway server is the place to do:

* Authentication
* Authorization
* Validations
* Throttling

Gateway server communicates with Monstera cluster via `monstera.MonsteraClient`. All Monstera operations are 
deterministic, so the gateway is the place to generate random numbers or get the current time __before__ sending a core
request to Monstera cluster.

This example is relatively simple and all operations from gateway API map 1-to-1 to core operations (not including 
authentication). However, in more complex applications, a single gateway operation can collect or update data in 
multiple application cores (Everblack Bison and Eveblack Moab has such operations, for example).

Here authentication is dumb. An account id is passed in headers and the server picks it up without any actual check. 
Just for demonstration purposes. In real Everblack Cloud it is implemented as described 
[here](https://everblack.dev/docs/api/authentication/), and 
[evrblk/evrblk-go/authn](https://github.com/evrblk/evrblk-go/tree/master/authn) package is used internally.

## Executables

The whole application consists of two parts:

* `cmd/gateway` - a stateless web server with public API. This is basically a runner for the Gateway gRPC server 
  from above.
* `cmd/node` - stateful Monstera node with all the data and business logic. This is a runner for 
  `monstera.MonsteraServer` and the place to register all implementations of your application cores.

Each Monstera node has two BadgerDB stores: one for all application cores, and one for all Raft logs from all shards
on that node.

`Procfile` has three Monstera nodes and one gateway server configured. Use goreman to start:

```
go tool github.com/mattn/goreman start
```

There is also a standalone executable `cmd/standalone` that runs both parts of the application in a single Go process,
non-sharded and non-replicated. Read more about standalone applications 
[here](https://everblack.dev/docs/monstera/standalone-application/). To run a standalone app:

```
go run ./cmd/standalone --port=8000 --data-dir=./data/standalone
```

## Monstera codegen

Monstera codegen is the opinionated part of the framework. I wanted to achieve type-safety and utilize compile-time 
checks without reflection, but wanted to eliminate human mistakes from vast boilerplate code. So I generate all 
boilerplate code.

`monstera.yaml` defines all application cores and their operations. `generator.go` has an annotation for running
`//go:generate` for Monstera codegen. It produces:

* `api.go` with interfaces for application cores and stubs
* `adapters.go` with adapter to application cores, that turns binary blobs into routable requests
* `stubs.go` with service stubs, that turn requests into binary blobs and route them to the correct application core

Monstera codegen relies on several conventions to make it work in a type-safe way:

* Methods of application core must have corresponding `*Response` and `*Request` objects in `go_code.corepb_package` 
  package. For example, `AcquireLock` of `Locks` core must have `AcquireLockRequest` and `AcquireLockResponse` proto 
  messages in `github.com/evrblk/monstera-example/dlocks/corepb`.
* `*Response` and `*Request` objects must be included into `oneof` of corresponding high level containers 
  `update_request_proto`, `update_response_proto`, etc. For example, `AcquireLockRequest` must be included into 
  `UpdateRequest`, `AcquireLockResponse` must be included into `UpdateResponse`.

The reason why I do not generate high level containers (in `corepb/cloud.proto`) is because of protobuf field tags.
They need to be consistent and never change. That means I would need to assign field tags right in the YAML file, which
I did not like. If I find an elegant and safe way to do it, I will simplify this codegen part.

`sharding.go` has an implementation of a shard key calculator. I chose not to use annotations or reflection to extract
shard keys from requests. Instead, Monstera codegen generates a simple interface where every method corresponds to 
a `*Request` object. You specify explicitly how to extract a shard key from each request with one line of Go code.

## Cluster config

Cluster config is used by MonsteraClient. There is already one generated for you in `cluster_config.pb`. 
`cluster_config.json` is a human-readable version of the same config, check it out.

To print a JSON version of any config run:

```
go tool github.com/evrblk/monstera/cmd/monstera cluster print-config --monstera-config=./cluster_config.pb
```

You can seed a new config. Keep in mind, if you run this command it will regenerate the config with new random ids, and 
you will also need to update `Procfile` with those new ids:

```
go run ./cmd/dev seed-monstera-cluster
```
