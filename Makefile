up: 
	docker compose up -d --build --remove-orphans

down:
	docker compose down

gate_way:
	go run gateway/*.go

user:
	go run service/user/*.go

video:
	go run service/video/*.go

relation:
	go run service/relation/*.go

message:
	go run service/message/*.go

favorite:
	go run service/favorite/*.go

comment:
	go run service/comment/*.go

proto:
	protoc --go_out=.. --go-grpc_out=.. ./idl/*.proto

model:
	go run cmd/main.go

tidy:
	go mod tidy && go fmt ./...