PG_DSN := "postgres://postgres:password@localhost:5452/vitalik?sslmode=disable"

.PHONY: build
build:
	@echo 'Building cmd/api...'
	go build -o=./bin/vitalik ./cmd/vitalik
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/vitalik ./cmd/vitalik

.PHONY: run
run:
	@go run ./cmd/vitalik

.PHONY: db-start
db-start:
	@echo 'Starting PostgreSQL database container using docker-compose...'
	docker-compose up -d db

.PHONY: db-stop
db-stop:
	@echo 'Stopping PostgreSQL database container...'
	docker-compose down

.PHONY: db-reset
db-reset: db-stop db-cleanup db-start wait-pg-up db-up

.PHONY: wait-pg-up
wait-pg-up:
	sleep 3

.PHONY: db-up
db-up:
	@echo 'Running database migrations...'
	goose -dir db/migrations postgres "$(PG_DSN)" up

.PHONY: db-cleanup
db-cleanup:
	@echo 'Cleaning up unused Docker volumes...'
	docker volume rm vitalik_backend_vitalik-db-data

.PHONE: db-jet
db-jet:
	@echo 'Generating go-jet code...'
	jet -source=postgres -host=localhost -port=5452 -user=postgres -password=password -dbname=vitalik -path=./.gen