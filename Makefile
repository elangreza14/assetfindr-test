#!make
include .env
	
run-http:
	go run cmd/http/main.go
	
stack-up:
	docker compose up -d

stack-down:
	docker compose down

.PHONY: run-http stack-up stack-down