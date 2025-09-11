include .env
export PATH := $(shell go env GOPATH)/bin:$(PATH)

.PHONY: migrate-up migrate-down

migrate-up:
	goose -dir mig/ postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir mig/ postgres "$(DATABASE_URL)" down

