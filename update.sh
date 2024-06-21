protoc --go_out=.. --go-grpc_out=.. ./idl/*.proto
go run cmd/main.go