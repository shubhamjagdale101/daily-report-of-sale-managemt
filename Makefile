.PHONY: run build test migrate-up migrate-down tidy

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

tidy:
	go mod tidy

docker-build:
	docker build -t gold-management-system .

docker-run:
	docker run -p 8080:8080 gold-management-system

migrate-up:
	# Add migration commands if using a migration tool

migrate-down:
	# Add migration commands if using a migration tool

rebuild:
	docker compose down --rmi local --remove-orphans
	docker compose build --no-cache
	docker compose up

rebuild-without-previous-volumes:
	docker compose down --volumes --rmi local --remove-orphans
	docker compose build --no-cache
	docker compose up