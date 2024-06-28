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

build:
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Version=$(version)" -mod=vendor -o stash stash.go

build-image:build
	docker build --rm --no-cache -t csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG) .

push-image:build-image
	docker push csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)
	echo "TKEx_Image_Hub=csighub.tencentyun.com/$(NAMESPACE)/scf-stash:$(TAG)" >> ${QCI_ENV_FILE}



