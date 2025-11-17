.PHONY: run build test migrate-up migrate-down clean

run:
	go run cmd/main.go

build:
	go build -o bin/server cmd/main.go

test:
	go test -v ./...

migrate-up:
	go run cmd/migrate/main.go -command up

migrate-down:
	go run cmd/migrate/main.go -command down

migrate-status:
	go run cmd/migrate/main.go -command status

migrate-to:
	go run cmd/migrate/main.go -command to -version $(VERSION)

clean:
	rm -rf bin/
	go clean
