#FROM golang:alpine AS builder
#
#LABEL stage=gobuilder
#
#ENV CGO_ENABLED 0
#ENV GOPROXY https://goproxy.cn,direct
#
#RUN apk update --no-cache && apk add --no-cache tzdata
#
#WORKDIR /build
#
#ADD go.mod .
#ADD go.sum .
#RUN go mod download
#COPY . .
#COPY stash/etc /app/etc
#RUN go build -ldflags="-s -w" -o /app/stash stash/stash.go
#
#
FROM centos:7

#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /app
COPY scf-stash /app/scf-stash
COPY etc /app/etc

CMD ["./scf-stash", "-f", "etc/config.yaml"]
