#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
protoc --go_out=. --go-grpc_out=. proto/kvstore.proto
echo "Generated proto Go files."
