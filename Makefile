build:
	@go build -o ./bin/ecommerce-psql ./cmd/main.go

test:
	@go test ./...

run: build
	@./bin/ecommerce-psql

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
