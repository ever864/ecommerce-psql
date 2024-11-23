postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=psql -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root ecommerce-psql

dropdb:
	docker exec -it postgres17 dropdb ecommerce-psql

migrate-down:
	migrate -path cmd/migrate/migrations -database "postgresql://root:psql@localhost:5432/ecommerce-psql?sslmode=disable" -verbose down

migrate-up:
	migrate -path cmd/migrate/migrations -database "postgresql://root:psql@localhost:5432/ecommerce-psql?sslmode=disable" -verbose up

build:
	@go build -o ./bin/ecommerce-psql ./cmd/main.go

test:
	@go test ./...

run: build
	@./bin/ecommerce-psql

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

.PHONY: createdb dropdb build test run migration migrate-up migrate-down postgres
