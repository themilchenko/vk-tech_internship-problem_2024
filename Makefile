ACTIVE_PACKAGES = $(go list ./... | grep -Ev "mocks" | tr '\n' ',')
MAIN_PATH = ./cmd/api/main.go
MOCKS_DESTINATION = internal/mocks

.PHONY: local_build
local_build: ## Build locally
	go build -o bin/ ${MAIN_PATH}

.PHONY: test
test: ## Run all the tests
	go test -v ./...

.PHONY: cover_out
cover_out: test ## Run all the tests and opens the coverage report
	go test -coverpkg=$(ACTIVE_PACKAGES) -coverprofile=c.out ./...
	cat c.out | grep -v "cmd" > tmp.out
	go tool cover -func=tmp.out

.PHONY: cover_html
cover_html: test ## Run all the tests and opens the coverage report in HTML
	go test -coverpkg=$(ACTIVE_PACKAGES) -coverprofile=c.out ./...
	cat c.out | grep -v "cmd" > tmp.out
	go tool cover -html=tmp.out

.PHONY: mocks
mocks: ## Generate mocks
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@mockgen -source=internal/domain/auth.go -destination=$(MOCKS_DESTINATION)/domain/auth.go
	@mockgen -source=internal/domain/actors.go -destination=$(MOCKS_DESTINATION)/domain/actors.go
	@mockgen -source=internal/domain/movies.go -destination=$(MOCKS_DESTINATION)/domain/movies.go
	@echo "OK"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := local_build
