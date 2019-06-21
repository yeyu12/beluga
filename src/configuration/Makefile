.PHONY:build
build:
	rm -rf $(shell pwd)/configuration
	$(GO_BUILD)
	mkdir -p $(BIN_DIR)/configuration
	cp -r $(shell pwd)/config $(BIN_DIR)/configuration/config
	cp -r $(shell pwd)/runtime $(BIN_DIR)/configuration/runtime
	cp -r $(shell pwd)/configuration $(BIN_DIR)/configuration/configuration
