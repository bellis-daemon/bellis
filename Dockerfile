# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.20 as builder

WORKDIR /workspace
# 启用go module
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ARG MODULE=backend
ARG PORT=7001

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -ldflags "-X 'main.GO_VERSION=$(go version)' -X 'main.BUILD_TIME=`TZ=Asia/Shanghai date "+%F %T"`'" -o entry modules/$MODULE/main.go


FROM ubuntu

EXPOSE $PORT

RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl
RUN DEBIAN_FRONTEND=noninteractive TZ=Asia/Shanghai apt-get -y install tzdata
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
#    && apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
ENV TZ=Asia/Shanghai
COPY --from=builder /go/bin/dlv /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert
COPY --from=builder /workspace/entry .
RUN chmod +x /entry

## 需要运行的命令
#ENTRYPOINT ["/dlv","--listen=:2345","--headless=true","--accept-multiclient","--api-version=2","exec","/entry"]
ENTRYPOINT ["/entry"]