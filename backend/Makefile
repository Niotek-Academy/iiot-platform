.PHONY: db-up db-down db-logs sqlc-generate migrate run tidy

# Start a local Postgres container matching .env.example defaults.
db-up:
	sudo docker run -d --name iiot-postgres \
		-e POSTGRES_USER=Niotek \
		-e POSTGRES_PASSWORD=Niotek \
		-e POSTGRES_DB=iiot_platform \
		-p 5432:5432 \
		postgres:12-alpine

db-down:
	sudo docker stop iiot-postgres && docker rm iiot-postgres

db-logs:
	sudo docker logs -f iiot-postgres

# Regenerate Go code from SQL queries + schema (requires sqlc installed).
sqlc-generate:
	sqlc generate

tidy:
	go mod tidy

run:
	go run ./cmd/server


DATABASE_URL ?= postgres://Niotek:Niotek@localhost:5432/iiot_platform?sslmode=disable

migrate-up:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" down 1

# Usage: make migrate-create name=add_new_column
migrate-create:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)
