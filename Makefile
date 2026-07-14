.PHONY: build test lint run dev clean docker-build docker-up

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

lint:
	golangci-lint run ./...

run:
	go run ./cmd/server

dev:
	air

clean:
	rm -rf bin/

docker-build:
	docker build -t open-payment-gateway .

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

migrate-up:
	go run ./internal/database/migrations/*.go up

migrate-down:
	go run ./internal/database/migrations/*.go down
