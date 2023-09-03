PROJECT_NAME := bias
CLI_NAME := $(PROJECT_NAME)-cli
API_NAME := $(PROJECT_NAME)-api
DIST_DIR := dist
CLI_SRC := cli/main.go
API_SRC := http/main.go
MAX_BINARIES ?= 3
BIAS_BUILD_CLI_CREATE_ALIAS ?= False

.PHONY: all clean list cli api

all: cli api

cli:
	@echo "Building $(CLI_NAME)..."
	@mkdir -p $(DIST_DIR)
	@go build -v -tags=exclude -o $(DIST_DIR)/$(CLI_NAME)-$(shell date +%s) $(CLI_SRC)
	@if [ "$(BIAS_BUILD_CLI_CREATE_ALIAS)" = "True" ]; then \
		cd $(DIST_DIR) && \
		ln -sf $$(ls -t $(CLI_NAME)-* | head -n 1) $(PROJECT_NAME); \
	fi

api:
	@echo "Building $(API_NAME)..."


clean:
	@echo "Cleaning old binaries..."
	@ls -t $(DIST_DIR)/$(CLI_NAME)-* | tail -n +$(shell echo $(MAX_BINARIES) + 1 | bc) | xargs rm -f
	@rm -f $(DIST_DIR)/$(PROJECT_NAME)
	@rm -f $(DIST_DIR)/$(PROJECT_NAME)-api

list:
	@ls -t $(DIST_DIR)/$(CLI_NAME)-*
