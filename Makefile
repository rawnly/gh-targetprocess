.PHONY: build install clean test run help

BINARY_NAME=gh-targetprocess
BUILD_DIR=.
GO=go

.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the extension binary
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)

install: build ## Install the extension locally using gh CLI
	gh extension install .

clean: ## Remove build artifacts
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

test: ## Run tests
	$(GO) test -v ./...

run: build ## Build and run the extension (use ARGS to pass arguments)
	./$(BINARY_NAME) $(ARGS)

uninstall: ## Uninstall the extension
	gh extension remove targetprocess