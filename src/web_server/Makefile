.PHONY:build
build:
	rm -rf $(shell pwd)/web_server
	$(GO_BUILD)
	mkdir -p $(BIN_DIR)/web_server
	cp -r $(shell pwd)/conf $(BIN_DIR)/web_server/conf
	cp -r $(shell pwd)/web_server $(BIN_DIR)/web_server/web_server
