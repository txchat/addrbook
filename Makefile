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

.PHONY: build build-arm64 clean

build:
	GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn" GOPRIVATE=gitlab.33.cn go build $(BUILD_FLAGS) -v -i -o $(APP) $(SRC_APP)
	GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn" GOPRIVATE=gitlab.33.cn go build $(BUILD_FLAGS) -v -i -o $(CLI) $(SRC_CLI)

build-arm64:
	GOOS=linux GOARCH=arm64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn" GOPRIVATE=gitlab.33.cn go build $(BUILD_FLAGS) -v -i -o $(APP)_linux_arm64 $(SRC_APP)
	GOOS=linux GOARCH=arm64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn" GOPRIVATE=gitlab.33.cn go build $(BUILD_FLAGS) -v -i -o $(CLI)_linux_arm64 $(SRC_CLI)

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