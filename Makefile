APP_NAME := static-serve-metadata
SOURCE_PATH := .

GOOS_VAR := linux
BIN_EXT :=

ifeq ($(OS), Windows_NT)
	GOOS_VAR := windows
	BIN_EXT := .exe
endif

# Some more documentation on this command for learning purpose:
# The `grep -E '^[a-zA-Z0-9_-]+:' Makefile`, this part finds any lines that matches as commands and its comments.
# For example: "help: ## Shows help command".
#
# The `awk` command has many instructions, so we'll split it:
# - `BEGIN { FS = ":( ##)?" };`, this sets the "file separator" to split the command and the comments.
# - `{ printf "\033[0;31m%-20s \033[0;32m%s\n", $$1, $$2 };`, this will print it as a nice looking help command.
.PHONY: help
help: ## Shows this command.
	@printf 'These are the available commands in our Makefile.\n'
	@printf '-------------------------------------------------\n'
	@grep -E '^[a-zA-Z0-9_-]+:' Makefile | awk 'BEGIN { FS = ":( ##)?" }; { printf "\033[0;31m%-20s \033[0m%s\n", $$1, $$2 };'

.PHONY: clean
clean: ## Cleans the build directory by removing all binary files.
	rm -rf build/*

.PHONY: build
build: ## Builds the app based on your operating system.
	go mod tidy -v
	GOOS=$(GOOS_VAR) go build -v -o ./build/$(APP_NAME)$(BIN_EXT) $(SOURCE_PATH)

.PHONY: build-prod
build-prod: ## Builds the app for production purpose.
	go mod tidy -v
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o ./build/$(APP_NAME)_linux_amd64 -v $(SOURCE_PATH)

.PHONY: build-all-prod
build-all-prod: ## Builds production-ready binaries for all platforms.
	@echo "Building production binaries for all platforms..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o ./build/$(APP_NAME)_linux_amd64 $(SOURCE_PATH)
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o ./build/$(APP_NAME)_linux_arm64 $(SOURCE_PATH)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o ./build/$(APP_NAME)_windows_amd64.exe $(SOURCE_PATH)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o ./build/$(APP_NAME)_darwin_amd64 $(SOURCE_PATH)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o ./build/$(APP_NAME)_darwin_arm64 $(SOURCE_PATH)

.PHONY: start
start: ## Starts the app from 'build' directory.
	ENVIRONMENT=production ./build/$(APP_NAME)$(BIN_EXT) start