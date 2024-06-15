.PHONY: build run_local run_tests run_docker run_docker_rebuild

build:
	go build ./cmd/api
run_local:
	go run ./cmd/api
run_tests:
	go test ./... -v
run_docker:
	docker-compose up -d
run_docker_rebuild:
	docker-compose up --build -d

