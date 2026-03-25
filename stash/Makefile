selfVer=1.0
commit=$(shell git rev-parse --short HEAD)
date=$(shell date +'%Y%m%d-%H:%M:%S')
gitversion=$(selfVer)-$(commit)-$(date)
goversion=$(shell go version | awk '{print $$3}')
version=$(gitversion)-$(goversion)

NAMESPACE=scf
COMMIT=$(shell git rev-parse --short HEAD)
DATE:=$(shell date +'%Y%m%d-%H%M%S')
TAG ?= release-stash-$(DATE)-$(COMMIT)
QCI_ENV_FILE?=release-tag

# 检测当前平台架构
NATIVE_ARCH := $(shell go env GOARCH)
NATIVE_OS := $(shell go env GOOS)

# 默认目标架构（可以通过命令行覆盖）
TARGETARCH ?= $(NATIVE_ARCH)
TARGETOS ?= linux

# 构建二进制文件（支持本地架构）
build:
	GOARCH=$(TARGETARCH) GOOS=$(TARGETOS) go build -ldflags "-X main.Version=$(version)" -mod=vendor -o stash stash.go

# 构建本地架构版本（用于本地测试）
build-native:
	GOARCH=$(NATIVE_ARCH) GOOS=$(NATIVE_OS) go build -ldflags "-X main.Version=$(version)" -mod=vendor -o stash stash.go

# 构建AMD64版本（用于Linux服务器部署）
build-amd64:
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Version=$(version)" -mod=vendor -o stash stash.go

# 构建ARM64版本
build-arm64:
	GOARCH=arm64 GOOS=linux go build -ldflags "-X main.Version=$(version)" -mod=vendor -o stash stash.go

# 使用Docker buildx进行多架构构建
build-image-multi:
	docker buildx build --platform linux/amd64,linux/arm64 --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG) .

# 传统单架构构建
build-image:build-amd64
	docker build --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG) .

# 本地开发用构建（匹配当前平台）
build-image-local:build-native
	docker build --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)-local .

# 推送单架构镜像
push-image:build-image
	docker push csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)
	echo "TKEx_Image_Hub=csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)" >> ${QCI_ENV_FILE}

# 设置多平台构建环境
setup-buildx:
	@echo "设置多平台构建环境..."
	docker buildx create --name multiarch --driver docker-container --use || docker buildx use multiarch
	docker buildx inspect --bootstrap

# 推送多架构镜像
push-image-multi: setup-buildx
	docker buildx build --platform linux/amd64,linux/arm64 --push --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG) .
	echo "TKEx_Image_Hub=csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)" >> ${QCI_ENV_FILE}

# 仅推送AMD64架构（适用于大多数Linux服务器）
push-image-amd64: build-amd64
	docker build --platform linux/amd64 --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG) .
	docker push csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)
	echo "TKEx_Image_Hub=csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)" >> ${QCI_ENV_FILE}

# 本地测试推送
push-image-local:build-image-local
	docker push csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)-local

# 清理构建产物
clean:
	rm -f stash
	docker rmi -f csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG) 2>/dev/null || true
	docker rmi -f csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)-local 2>/dev/null || true

# 显示帮助信息
help:
	@echo "可用的构建目标："
	@echo "  build              - 构建二进制文件（默认Linux，架构可通过TARGETARCH指定）"
	@echo "  build-native       - 构建本地平台的二进制文件"
	@echo "  build-amd64        - 构建AMD64架构的Linux二进制文件"
	@echo "  build-arm64        - 构建ARM64架构的Linux二进制文件"
	@echo "  build-image        - 构建Docker镜像（AMD64）"
	@echo "  build-image-local  - 构建本地平台的Docker镜像"
	@echo "  build-image-multi  - 构建多架构Docker镜像"
	@echo "  push-image         - 推送Docker镜像（原版，AMD64）"
	@echo "  push-image-amd64   - 推送AMD64架构Docker镜像（推荐用于Linux服务器）"
	@echo "  push-image-multi   - 推送多架构Docker镜像（自动设置buildx环境）"
	@echo "  setup-buildx       - 设置多平台构建环境"
	@echo "  clean              - 清理构建产物"
	@echo ""
	@echo "🚀 推荐用法："
	@echo "  Linux服务器部署：make push-image-amd64"
	@echo "  本地开发测试：  make build-native && make build-image-local"
	@echo "  多平台发布：    make push-image-multi"
	@echo ""
	@echo "环境变量："
	@echo "  TARGETARCH=amd64|arm64  - 目标架构（默认：$(NATIVE_ARCH)）"
	@echo "  TARGETOS=linux|darwin   - 目标操作系统（默认：linux）"
	@echo ""
	@echo "当前平台：$(NATIVE_OS)/$(NATIVE_ARCH)"

.PHONY: build build-native build-amd64 build-arm64 build-image build-image-local build-image-multi push-image push-image-amd64 push-image-multi push-image-local setup-buildx clean help

.DEFAULT_GOAL := help
