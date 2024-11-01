include .env
export $(shell sed 's/=.*//' .env)

MIGRATIONS_DIR = ./schema

DB_CONN := "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DATABASE} sslmode=${POSTGRES_SSL_MODE}"

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONN) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONN) down

swag:
	swag init -g cmd/main.go

start:
	docker-compose up -d

stop:
	docker-compose down

build:
	docker-compose build