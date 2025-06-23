package dlocks

var (
	// Unique table prefixes to avoid collisions between tables while sharding the same instance of BadgerDB store.
	accountsTableId       = []byte{0x00, 0x00}
	accountsEmailsIndexId = []byte{0x00, 0x01}

	locksTableId = []byte{0x01, 0x00}

	namespacesTableId = []byte{0x02, 0x00}
)
