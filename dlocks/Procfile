node01: go run ./cmd/node --port=7000 --node-id=nd_c95f3e0f --data-dir=./data/nd_c95f3e0f --monstera-config=./cluster_config.pb
node02: go run ./cmd/node --port=7001 --node-id=nd_edbe72ff --data-dir=./data/nd_edbe72ff --monstera-config=./cluster_config.pb
node03: go run ./cmd/node --port=7002 --node-id=nd_d64722af --data-dir=./data/nd_d64722af --monstera-config=./cluster_config.pb

gateway: go run ./cmd/gateway --port=8000 --monstera-config=./cluster_config.pb
