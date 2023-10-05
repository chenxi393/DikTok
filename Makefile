run:
	go run .

up: 
	docker-compose up -d --build --remove-orphans

down:
	docker-compose down

test: 
	go run main.go route.go

