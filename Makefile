.PHONY: help build test coverage mocks mocks-install docker-up docker-down run work-init work-sync work-tidy work-vendor

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  test           Run unit tests"
	@echo "  coverage       Run tests with coverage report"
	@echo "  mocks          Generate mocks using mockgen"
	@echo "  mocks-install  Install mockgen tool"
	@echo "  work-init      Initialize go.work file with all modules"
	@echo "  work-sync      Sync dependencies across all modules"
	@echo "  work-tidy      Tidy go.mod files in all modules"
	@echo "  work-vendor    Vendor dependencies for all modules"

build:
	go build -o bin/api ./cmd/api

test:
	go test -v -race -timeout 30s ./domains/...
	go test -v -race -timeout 30s ./ports/...
	go test -v -race -timeout 30s ./services/...

coverage:
	go test -v -race -timeout 30s -coverprofile=coverage.out ./domains/...
	go test -v -race -timeout 30s -coverprofile=coverage.out ./ports/... -coverappend
	go test -v -race -timeout 30s -coverprofile=coverage.out ./services/... -coverappend
	go tool cover -html=coverage.out -o coverage.html

clear-mocks:
	rm -rf testsuites/mocks/*

mocks: mocks-install clear-mocks
	mkdir -p testsuites/mocks
	@find ports -name "*.go" -type f | while read file; do \
		filename=$$(basename $$file .go); \
		mockgen -source=$$file -destination=testsuites/mocks/mock_$$filename.go -package=mocks; \
	done

work-init:
	go work init
	go work use ./domains ./ports ./services ./adapters ./testsuites

work-sync:
	go work sync

work-tidy:
	go mod tidy -C ./domains
	go mod tidy -C ./ports
	go mod tidy -C ./services
	go mod tidy -C ./adapters
	go mod tidy -C ./testsuites

work-vendor:
	go mod vendor -C ./domains
	go mod vendor -C ./ports
	go mod vendor -C ./services
	go mod vendor -C ./adapters
	go mod vendor -C ./testsuites

.DEFAULT_GOAL := help
