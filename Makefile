.PHONY: help build test coverage mocks mocks-install docker-up docker-down run

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build          Build the API"
	@echo "  test           Run unit tests"
	@echo "  coverage       Run tests with coverage report"
	@echo "  mocks          Generate mocks using mockgen"
	@echo "  mocks-install  Install mockgen tool"
	@echo "  docker-up      Start Docker Compose services"
	@echo "  docker-down    Stop Docker Compose services"
	@echo "  run            Run the API"

build:
	go build -o bin/api ./cmd/api

test:
	go test -v -race -timeout 30s ./...

coverage:
	go test -v -race -timeout 30s -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

mocks-install:
	go install github.com/golang/mock/mockgen@latest

clear-mocks-old:
	rm -rf mocks/*

mocks-old: mocks-install clear-mocks-old
	mkdir -p mocks
	@find internal/domain -name "*.go" -type f | grep -E "(service|repository)" | while read file; do \
		filename=$$(basename $$file .go); \
		mockgen -source=$$file -destination=mocks/mock_$$filename.go -package=mocks; \
	done

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

run: build
	./bin/api

clear-mocks:
	rm -rf testsuites/mocks/*

mocks: mocks-install clear-mocks
	mkdir -p testsuites/mocks
	@find ports -name "*.go" -type f | while read file; do \
		filename=$$(basename $$file .go); \
		mockgen -source=$$file -destination=testsuites/mocks/mock_$$filename.go -package=mocks; \
	done

.DEFAULT_GOAL := help
