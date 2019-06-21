PACKAGE_LOCK = package-lock.json

.PHONY:build
build:
	rm -rf dist/*

	@echo "##############\n# 安装依赖   #\n##############\n"
	npm install
	@echo "##############\n# 开始编译   #\n##############\n"
	npm run build

	# 文件copy到build目录中去
	mv $(shell pwd)/dist $(BIN_DIR)/web_ui

.PHONY:clear
clear:
	rm -rf dist/*
	rm -rf node_modules/
	rm -rf package-lock.json
