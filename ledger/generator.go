package ledger

// Generate protos
//go:generate bash ./gen_protos.sh

// Generate Monstera stubs, APIs, and adapters
//go:generate go tool github.com/evrblk/monstera/cmd/monstera code generate
