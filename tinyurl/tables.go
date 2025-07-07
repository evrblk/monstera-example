package tinyurl

var (
	// Unique table prefixes to avoid collisions between tables while sharding the same instance of BadgerDB store.
	usersTableId     = []byte{0x00, 0x00}
	shortUrlsTableId = []byte{0x01, 0x00}
)
