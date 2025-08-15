.PHONY: proto tidy run1 run2 run3

proto:
	protoc --go_out=. --go-grpc_out=. proto/kvstore.proto

tidy:
	go mod tidy

run1:
	go run server/main.go --node_id=node1 --port=50051
run2:
	go run server/main.go --node_id=node2 --port=50052
run3:
	go run server/main.go --node_id=node3 --port=50053
