.PHONY:build
build:
	rm -rf $(shell pwd)/task
	$(GO_BUILD)
	mkdir -p $(BIN_DIR)/task
	cp -r $(shell pwd)/config $(BIN_DIR)/task/config
	cp -r $(shell pwd)/task $(BIN_DIR)/task/task
