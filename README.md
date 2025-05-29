# Monstera Example

An example of how to build applications with Monstera framework. This is an imaginary multi-tenant SaaS for 
distributed RW locks. Locks are referenced by a user-specified name and are organized into namespaces. Names for locks
are unique within a namespace, but not between namespaces. Basically, this is a simplified version of Everblack Grackle
service, with locks only, trivial account management, and no authentication. 

Monstera framework does not force any particular application core implementation, method routing mechanism, or any 
specific wire format. It is up to you to define that. However, over time I developed a certain style of how all
Everblack services are implemented. To separate a clean part of the framework from that opinionated part I made
two packages: `github.com/evrblk/monstera` for the core part, and `github.com/evrblk/monstera/x` for the rest.
However, a lot of things are not generalizable or extractable into a library. And this example application shows how
all of them can be assembled together.

## Applications cores

There are 3 application cores:

* `AccountsCore` in `./accounts.go`
* `NamespacesCore` in `./namespaces.go`
* `LocksCore` in `./locks.go`

Application cores store data in BadgerDB. There is one instance of BadgeDB per process, so multiple shards and multiple
cores share it. To avoid conflicts, each table is prefixed with table IDs (in `./tables.go`). Each shard has its own 
boundaries (`lowerBound` and `upperBound`). Take a look how keys are built for tables and indexes (typically in the
bottom of files with application cores and inside `monstera/x` package too).

All core data structures are defined in protobufs in `./corepb/*`. Those structures are exposed from Monstera stubs
and used by application cores to store data in BadgerDB. `./corepb/cloud.proto` has high level containers for requests
and responses that are actually passed by Monstera. Monstera does not know anything about implementation of your
application cores and only passes binary blobs as requests and responses for reads and updates.

## Gateway server

A gateway (or frontend) server is the public API part of the system. In this example it serves gRPC, but it can be
anything you want (OpenAPI, ConnectRPC, gRPC, Gin, etc). Gateway gRPC is defined in `./gatewaypb/`. Protos are not
shared between gateway and core parts for clean separation of core business layer and presentation layer. The code for
converting between them lives in `./server/pbconv.go`. `./server/server.go` is the implementation of the gateway API.
It is the entry point for all user actions, and if you want to trace and understand the lifecycle of a request then 
start from here.

Gateway server is the place to do:

* Authentication
* Authorization (not in this example)
* Validations
* Throttling (not in this example)

Gateway server communicates with Monstera cluster via `monstera.MonsteraClient`. All Monstera operations are 
deterministic, so the gateway is the place to generate random numbers or get the current time __before__ sending a core
request to Monstera cluster.

This example is relatively simple and all operations from gateway API map 1-to-1 to core operations (not including 
authentication). However, in more complex applications a single gateway operation can collect or update data in 
multiple application cores (Everblack Bison and Eveblack Moab has such operations, for example).

## Executables

The whole application consists of two executables:

* `./cmd/gateway` - a stateless web server with public API. This is basically a runner for the Gateway gRPC server 
  from above.
* `./cmd/node` - stateful Monstera node with all the data and business logic. This is a runner for 
  `monstera.MonsteraServer` and the place to register all implementations of your application cores.
