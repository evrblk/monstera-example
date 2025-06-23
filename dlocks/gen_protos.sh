# Generate ./corepb protos
# Protoc should know where to find shared types from github.com/evrblk/monstera
MONSTERA_PROTO_ROOT=$(go list -f '{{.Dir}}' -m github.com/evrblk/monstera)
protoc \
    --proto_path . \
    --proto_path "$MONSTERA_PROTO_ROOT" \
    --go_out=. \
    --go_opt=paths=source_relative \
    ./corepb/*.proto

# Generate ./gatewaypb protos
protoc \
    --proto_path . \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    ./gatewaypb/*.proto
