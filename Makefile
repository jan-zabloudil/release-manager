.PHONY: build run_local

build:
	go build ./cmd/api
run_local:
	go run ./cmd/api
