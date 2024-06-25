#!make
include .env
	
run-http:
	go run cmd/http/main.go
	
stack-up:
	docker compose up -d

stack-down:
	docker compose down

gen:
	go generate ./...

test-cover:
	go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out

.PHONY: run-http stack-up stack-down gen test-coverage