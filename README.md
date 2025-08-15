# Tenac – Week 2 Drop-in

This bundle includes:
- Updated `proto/kvstore.proto` with `HealthCheck`
- Fixed `server/main.go` using modern gRPC creds
- `internal/store.go`
- `client/main.go`
- `nodes.json`
- `scripts/gen.sh` and `Makefile` to generate proto code

## Quick start

1) Install protoc + plugins (once):
   ```bash
   brew install protobuf   # macOS (or download from GitHub releases)
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

2) Generate gRPC code:
   ```bash
   make proto
   # or
   ./scripts/gen.sh
   ```

3) Run 3 nodes in separate terminals:
   ```bash
   go run server/main.go --node_id=node1 --port=50051
   go run server/main.go --node_id=node2 --port=50052
   go run server/main.go --node_id=node3 --port=50053
   ```

4) (Optional) Test client:
   ```bash
   go run client/main.go --addr localhost:50051
   ```
