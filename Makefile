export GO111MODULE=on

VERSION := $(shell echo $(shell cat version/version.go | grep "Version =" | cut -d '=' -f2))
BUILD_DIR := build
CLI := ${BUILD_DIR}/chain33-cli
SRC_CLI := github.com/txchat/addrbook/cmd/cli
APP := ${BUILD_DIR}/chain33
SRC_APP:= github.com/txchat/addrbook/cmd/chain33
BUILD_FLAGS = -ldflags "-X github.com/33cn/chain33/common/version.GitCommit=`git rev-parse --short=8 HEAD`"
LDFLAGS := -ldflags "-w -s"
APP_NAME := chain33
PKG_NAME := ${APP_NAME}_v${VERSION}
PKG := ${PKG_NAME}.tar.gz

HOSTARCH=$(shell uname -m)
HOSTOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOBASE=GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn"
GOENV=GOOS=$(HOSTOS) GOARCH=$(HOSTARCH) ${GOBASE}


.PHONY: build build-arm64 clean

help: ## Display this help screen
	@printf "Help doc:\nUsage: make [command]\n"
	@printf "[command]\n"
	@grep -h -E '^([a-zA-Z_-]|\%)+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:	## 编译本机系统和指令集的可执行文件
	$(GOENV) go build $(BUILD_FLAGS) -v -i -o $(APP) $(SRC_APP)
	$(GOENV) go build $(BUILD_FLAGS) -v -i -o $(CLI) $(SRC_CLI)

build_%:  ## 编译目标机器的可执行文件（例如: make build_linux_amd64）
	TAR_OS=$(shell echo $* | awk -F'_' '{print $$1}'); \
	TAR_ARCH=$(shell echo $* | awk -F'_' '{print $$2}'); \
	GOOS=$${TAR_OS} GOARCH=$${TAR_ARCH} $(GOBASE) go build $(BUILD_FLAGS) -v -i -o $(APP) $(SRC_APP);\
	GOOS=$${TAR_OS} GOARCH=$${TAR_ARCH} $(GOBASE) go build $(BUILD_FLAGS) -v -i -o $(CLI) $(SRC_CLI)

pkg: build
	mkdir -p ${PKG_NAME}
	cp ${BUILD_DIR}/chain33 ${PKG_NAME}/
	cp ${BUILD_DIR}/chain33-cli ${PKG_NAME}/
	cp ${BUILD_DIR}/*.toml ${PKG_NAME}/
	tar zvcf ${PKG} ${PKG_NAME}
	rm -rf ${PKG_NAME}

clean:
	@rm -rf build/datadir
	@rm -rf build/chain33*
	@rm -rf build/*.log
	@rm -rf build/*.toml
	@rm -rf build/logs
	@rm -rf build/wallet
	@go clean