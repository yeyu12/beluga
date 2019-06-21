GO_GET = go get -v -u
GO_RELY_ON = dep ensure -v
GO_BUILD = go build
NPM_BUILD = npm run build
PWD_DIR = $(PWD)
BIN_DIR = $(PWD_DIR)/bin
SUB_DIR = src/web_ui \
			src/web_server \
			src/configuration/

GOLANG_X = $(GOPATH)/src/golang.org/x
GOLANG_GITHUB = $(GOPATH)/src/github.com
GOLANG_ETCD = $(GOLANG_GITHUB)/coreos
GOLANG_GRPC = $(GOPATH)/src/google.golang.org/
GOLANG_PROJECT_GRPC = $(PWD_DIR)/vendor/google.golang.org
GIT_CLONE = git clone

export GO_GET GO_RELY_ON GO_BUILD NPM_BUILD PWD_DIR BIN_DIR SUB_DIR

all: $(SUB_DIR)

$(SUB_DIR):ECHO
	@echo "编译$@"
	make -C $@

ECHO:
	@echo $(SUB_DIR)
	@echo begin compile

	# 删除bin目录下的所有文件
	rm -rf $(BIN_DIR)/*

	# 检测环境
	# 安装go依赖
	$(GO_GET) github.com/golang/dep/cmd/dep

ifeq "$(wildcard $(BIN_DIR))" ""
	mkdir -p $(BIN_DIR)
endif

ifeq "$(wildcard $(GOLANG_X))" ""
	# 安装go官方的包
	mkdir -p $(GOLANG_X)

	# golang官方包
	$(GIT_CLONE) https://github.com/golang/sys.git $(GOLANG_X)/sys
	$(GIT_CLONE) https://github.com/golang/crypto.git $(GOLANG_X)/crypto
	$(GIT_CLONE) https://github.com/golang/net.git $(GOLANG_X)/net
	$(GIT_CLONE) https://github.com/golang/image.git $(GOLANG_X)/image
	$(GIT_CLONE) https://github.com/golang/text.git $(GOLANG_X)/text
endif

# 判断etcd是否在gopath目录下面。先go get一下，直接dep会存在依赖包的问题
ifeq "$(wildcard $(GOLANG_ETCD))" ""
	$(GIT_CLONE) https://github.com/etcd-io/etcd.git $(GOLANG_GITHUB)/coreos/etcd
endif

	$(GO_RELY_ON)

# 判断grpc是否在vendor目录下面,这里要整理下.
ifeq "$(wildcard $(GOLANG_PROJECT_GRPC))" ""
ifeq "$(wildcard $(GOPATH)/src/google.golang.org/grpc)" ""
	$(GIT_CLONE) https://github.com/grpc/grpc-go.git $(GOPATH)/src/google.golang.org/grpc
endif
ifeq "$(wildcard $(GOPATH)/src/google.golang.org/genproto)" ""
	$(GIT_CLONE) https://github.com/google/go-genproto.git $(GOPATH)/src/google.golang.org/genproto
endif

	$(GO_GET) github.com/golang/protobuf/ptypes

	cp -rf $(GOLANG_GRPC) $(GOLANG_PROJECT_GRPC)
endif



