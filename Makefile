# ==============================
# Makefile untuk Project Go
# ==============================

# Cari semua folder yang punya file *_test.go
TEST_PKGS := $(shell find . -name '*_test.go' -exec dirname {} \; | sort -u)

# ==============================
# PHONY targets
# ==============================
.PHONY: run seed build test cover test-report migrate-up migrate-down migrate-status migrate-to clean

# ==============================
# Run aplikasi
# ==============================
run:
	@echo "Running main application..."
	go run cmd/main.go

# ==============================
# Seed database
# ==============================
seed:
	@echo "Seeding database..."
	go run cmd/seed/main.go

# ==============================
# Build aplikasi
# ==============================
build:
	@echo "Building application..."
	go build -o bin/server cmd/main.go

# ==============================
# Jalankan test hanya package dengan *_test.go
# ==============================
test:
	@echo "Running tests for packages with *_test.go only..."
	@go test -v $(TEST_PKGS)

# ==============================
# Jalankan test dengan laporan JSON (untuk CI)
# ==============================
test-report:
	@echo "Running tests with JSON output..."
	@go test -v ./... -json | go-test-report

# ==============================
# Jalankan test + coverage
# ==============================
cover:
	@echo "Running coverage for packages with *_test.go only..."
	@go test -coverprofile=coverage.out $(TEST_PKGS)
	@echo "Opening HTML coverage report..."
	@go tool cover -html=coverage.out

# ==============================
# Database migration commands
# ==============================
migrate-up:
	@echo "Running migration up..."
	go run cmd/migrate/main.go -command up

migrate-down:
	@echo "Running migration down..."
	go run cmd/migrate/main.go -command down

migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrate/main.go -command status

migrate-to:
	@echo "Migrating to version $(VERSION)..."
	go run cmd/migrate/main.go -command to -version $(VERSION)

# ==============================
# Clean project
# ==============================
clean:
	@echo "Cleaning project..."
	rm -rf bin/
	go clean


# .PHONY: run build test migrate-up migrate-down clean

# run:
# 	go run cmd/main.go
	
# seed:
# 	go run cmd/seed/main.go

# build:
# 	go build -o bin/server cmd/main.go

# test:
# 	go test -v ./...

# test-report:
# 	go test -v ./... -json | go-test-report

# migrate-up:
# 	go run cmd/migrate/main.go -command up

# migrate-down:
# 	go run cmd/migrate/main.go -command down

# migrate-status:
# 	go run cmd/migrate/main.go -command status

# migrate-to:
# 	go run cmd/migrate/main.go -command to -version $(VERSION)

# clean:
# 	rm -rf bin/
# 	go clean
