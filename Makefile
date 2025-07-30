.PHONY: help
help:
	@echo "🛠️ Dev Commands\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

BIN_DIR=$(PWD)/tmp/bin
GOBIN ?= $(BIN_DIR)

export GOBIN
export PATH := $(BIN_DIR):$(PATH)

.PHONY: tools
tools: tools-test ## Install static checkers & other binaries
	@echo "🚚 Downloading tools.."
	@mkdir -p $(GOBIN)
	@ \
	command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest & \
	command -v goreleaser > /dev/null || go install github.com/goreleaser/goreleaser/v2@latest & \
	wait

.PHONY: lint
lint: tools ## Lint the source code
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 GoCI Lint.."
	@$(BIN_DIR)/golangci-lint fmt ./...
	@$(BIN_DIR)/golangci-lint run ./...
	@echo "🧹 Check GoReleaser.."
	@$(BIN_DIR)/goreleaser check

.PHONY: lint-ci
lint-ci: ## Lint the source code in CI mode
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@echo "🧹 Vetting go.mod.."
	@go vet ./...

.PHONY: test
test: ## Run tests
	@go test -v -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...