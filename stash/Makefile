APP = scf-stash

OS = $(shell uname -s)
ifeq ($(OS),Linux)
	GO_ENV=CGO_ENABLED=1 GOOS=linux GOARCH=amd64
else ifeq ($(OS),Darwin)
	# use "brew install FiloSottile/musl-cross/musl-cross" install linux-musl-gcc
	GO_ENV=CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" GOOS=linux GOARCH=amd64
endif

NAMESPACE=scf
COMMIT=$(shell git rev-parse --short HEAD)
DATE:=$(shell date +'%Y%m%d-%H%M%S')
QCI_ENV_FILE?=release-tag

TAG ?= release-data-update-$(DATE)-$(COMMIT)

build-docker:
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Version=$(version)" -mod=vendor -o $(APP) stash.go

build-image:build-docker
	docker build --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/$(APP):$(TAG) .

push-image:build-image
	docker push csighub.tencentyun.com/$(NAMESPACE)/$(IMG-RegionInvoke):$(TAG)
	echo "TKEx_Image_Hub=csighub.tencentyun.com/$(NAMESPACE)/$(IMG-RegionInvoke):$(TAG)" >> ${QCI_ENV_FILE}

