GO := go
BINARY_NAME := server
BUILD_DIR := build

build:
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

dev:
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .
	$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@rm -rf $(BUILD_DIR)
clean-all:
	@rm -rf $(BUILD_DIR)

.PHONY: build dev clean clean-all
