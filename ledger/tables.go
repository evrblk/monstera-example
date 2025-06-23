package ledger

var (
	// Unique table prefixes to avoid collisions between tables while sharding the same instance of BadgerDB store.
	accountsTableId     = []byte{0x00, 0x00}
	transactionsTableId = []byte{0x00, 0x01}
)
