GO_CMD = go
APP_NAME = gopod
BIN_DIR = bin

build:
	@echo "Building $(APP_NAME)..."
	$(GO_CMD) build -o $(BIN_DIR)/$(APP_NAME) .

clean:
	@echo "Cleaning up"
	$(GO_CMD) clean
	rm -rf $(BIN_DIR)

.PHONY: all build clean
