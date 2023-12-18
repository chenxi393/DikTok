run:
	go run .

up: 
	docker compose up -d --build --remove-orphans

down:
	docker compose down

gate_way:
	go run gateway/main.go gateway/msg.go gateway/route.go

user:
	go run service/user/main.go service/user/user.go

video:
	go run service/video/main.go service/video/publish.go service/video/video.go service/video/search.go

relation:
	go run service/relation/main.go service/relation/relation.go

message:
	go run service/message/main.go service/message/message.go

favorite:
	go run service/favorite/main.go service/favorite/favorite.go

comment:
	go run service/comment/main.go service/comment/comment.go 