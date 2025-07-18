# TinyURL Service

[![Go Report Card](https://goreportcard.com/badge/github.com/evrblk/monstera-example/tinyurl)](https://goreportcard.com/report/github.com/evrblk/monstera-example/tinyurl)

URL shortener service with multiple users. Users can list their URLs. Short URLs are base62 encoded 8 bytes. This simple
version has 4 bytes for a user ID and 4 bytes for a URL ID. It can be tuned up or down based on expected number of users
or URLs per user. If there were no requirement to be able to list URLs for a user then user id could be removed 
completely from the system. Here short URLs have user ID inside them so it is possible to route `GetShortUrl` requests 
to corresponding shards based on that user ID. Also `GetShortUrl` reads from followers (eventually consistent) to 
increase read throughput of the system.

![Diagram](diagram.png)

## Application cores

There are two application cores:

* `Users` in `users.go`. Sharded by user id.
  * `GetUser`
* `ShortUrls` in `shorturls.go`. Sharded by user id.
  * `GetShortUrl` (reads from followers are allowed)
  * `ListShortUrls`
  * `CreateShortUrl`

Take a look at tests (`users_test.go`, `shorturls_test.go`). 

## Cluster config

Pregenerated cluster config `cluster_config.json` has:

* 3 nodes
* 16 shards of `ShortUrls`
* 4 shards of `Users`
* 3 replicas of each

## How to run

1. Clone this repository.

```
git clone git@github.com:evrblk/monstera-example.git

cd ./monstera-example/tinyurl
```

2. Make sure it builds:

```
go build -v ./...
```

3. Start a cluster with 3 nodes and a gateway server:

```
go tool github.com/mattn/goreman start
```

4. Create 100 users:

```
go run ./cmd/dev seed-users
```

5. Pick any user id from the previous step output.

6. Run test scenario 1 which creates a short URL for the user id:

```
go run ./cmd/dev scenario-1 --user-id=9fff3bf7d1f9561d
```

## How to explore

For example, you want to understand how `CreateShortUrl` method works:

* Start reading from `server.go` in `TinyUrlServiceApiServer.CreateShortUrl()`.
* Trace it down to `monstera.MonsteraClient` calls.
* Optional: You can jump further if you want to read Monstera sources.
* Find `ShortUrlsCoreAdapter.Update()` and find `CreateShortUrl` there.
* Trace it down to `ShortUrlsCore.CreateShortUrl()`.
* Understand how simple the code is and how it takes advantage of sequential application of updates.
