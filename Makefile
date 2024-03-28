run:
	go run .

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
	go build service/message/*.go && ./biz

favorite:
	go run service/favorite/*.go

comment:
	go run service/comment/*.go