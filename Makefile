PROJECT_NAME := "go-types"
PKG := "github.com/a-castellano/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all build clean test test_integration test_rabbitmq coverage coverhtml lint race msan

all: build

lint: ## Lint the files
	@go vet ${PKG_LIST}

test: ## Run unit tests
	@go test --tags=unit_tests -short ./...

test_integration: ## Run integration tests
	@go test --tags=integration_tests -short ./...

test_rabbitmq: ## Run rabbitmq related tests
	@go test --tags=rabbitmq_tests -short ./...

test_rabbitmq_unit: ## Run rabbitmq related tests
	@go test --tags=rabbitmq_unit_tests -short ./...

test_redis: ## Run redis related tests
	@go test --tags=redis_tests -short ./...

test_redis_unit: ## Run redis related tests
	@go test --tags=redis_unit_tests -short ./...

race: ## Run data race detector
	@go test -race -short ./...

msan: ## Run memory sanitizer
	@go test -msan -short ./...

coverage: ## Generate global code coverage report
	./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./scripts/coverage.sh html;

#build: ## Build the binary file
#	@go build -v $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

