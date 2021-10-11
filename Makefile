compile-protoc-go:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
run-server:
	go run cmd/server/server.go
run-client:
	go run cmd/client/client.go