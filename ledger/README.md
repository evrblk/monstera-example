# Ledger Service

Simple bookkeeping service, it has multiple accounts, tracks their balances and allows to create transactions.

A transaction can be created with `CreateTransaction`. Transaction `amount` will be added to the account balance.
Positive transactions are debits (or "topups"), negative transactions are credits (or "purchases"). A transaction can be
instant (with `settled = true`) or pending (with `settled = false`). Pending transactions will be settled or canceled 
later with `SettleTransaction` and `CancelTransaction` respectively.

Account balance is represented with `available_balance` and `settled_balance` fields. Settled transaction affect both
balances. Pending credits are subtracted from the avaialble balance only. Each negative transaction is checked against 
available balance and `INSUFFICIED_FUNDS` error will be returned if the balance goes negative after such a transaction.
It is possible to topup negative balance to any amount.

Compared to other popular approaches to solve Ledger System Design interview questions this approach:

* has realtime account balance (it is updated instantly after each transaction is processed)
* allows for insufficient funds check (precisely and instantly)
* allows for complex logic around available and settled balances, or around negative balances 
* has fewer moving parts (no streams, no async workers)
* scales infinitely by the number of accounts


## Application cores

There is 1 application core:

* `AccountsCore` in `accounts.go`

Take a look at tests (`accounts_test.go`). 

## Cluster config

Pregenerated cluster config `cluster_config.json` has:

* 3 nodes
* 16 shards of `Accounts`
* 3 replicas of each

## How to run

1. Clone this repository.

```
git clone git@github.com:evrblk/monstera-example.git

cd ./monstera-example/ledger
```

2. Make sure it builds:

```
go build -v ./...
```

3. Start a cluster with 3 nodes and a gateway server:

```
go tool github.com/mattn/goreman start
```

4. Create 100 accounts:

```
go run ./cmd/dev seed-accounts
```

5. Pick any account id from previous step output.

6. Run a test scenario 1 which creates a namespace and tries to grab a lock with the account id:

```
go run ./cmd/dev scenario-1 --account-id=9fff3bf7d1f9561d
```

## How to explore

For example, you want to understand how `CreateTransaction` method works:

* Start reading from `server.go` in `LedgerServiceApiServer.CreateTransaction()`.
* Trace it down to `monstera.MonsteraClient` calls.
* Optional: You can jump further if you want to read Monstera sources.
* Find `AccountsCoreAdapter.Update()` and find `CreateTransaction` there.
* Trace it down to `AccountsCore.CreateTransaction()`.
* Understand how simple the code is and how it takes advantage of sequential application of updates.
