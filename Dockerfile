# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.21.0-alpine as builder

WORKDIR /workspace

RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*

# 启用go module
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct

# RUN go install github.com/go-delve/delve/cmd/dlv@latest

ARG MODULE=backend
ARG PORT=7001

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -ldflags "-s -w -X 'main.GO_VERSION=$(go version)' -X 'main.BUILD_TIME=`TZ=Asia/Shanghai date "+%F %T"`'" -o entry modules/$MODULE/main.go \
    && upx -9 entry

FROM alpine

EXPOSE $PORT
ENV TZ=Asia/Shanghai
#RUN DEBIAN_FRONTEND=noninteractive TZ=Asia/Shanghai apt-get -qq update \
#    && apt-get -qq install -y --no-install-recommends ca-certificates curl tzdata
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++ curl
#COPY --from=builder /go/bin/dlv /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert
COPY --from=builder --chmod=777 /workspace/entry .
COPY --chmod=777 hostname.sh .
## 需要运行的命令
#ENTRYPOINT ["/dlv","--listen=:2345","--headless=true","--accept-multiclient","--api-version=2","exec","/entry"]
ENTRYPOINT ["/entry"]