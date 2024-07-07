FROM golang:1.22.5-alpine as builder

WORKDIR /workspace

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories 
RUN apk add --no-cache upx tzdata && rm -rf /var/cache/apk/*

ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn
ENV TZ=Asia/Shanghai

ARG MODULE=backend

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -gcflags "all=-N -l" -ldflags "-s -w -X 'main.GoVersion=$(go version)' -X 'main.BuildTime=$(date "+%F %T")'" -o entry modules/$MODULE/main.go \
    && upx -6 entry

FROM envoyproxy/envoy:distroless-v1.30.4 AS envoy


FROM chromedp/headless-shell:128.0.6559.0
EXPOSE 7001
ENV TZ=Asia/Shanghai
COPY --from=envoy --chmod=755 /usr/local/bin/envoy /usr/local/bin
RUN DEBIAN_FRONTEND=noninteractive TZ=Asia/Shanghai apt-get -qq update \
   && apt-get -qq install -y --no-install-recommends ca-certificates curl tzdata 
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
#     && apk upgrade --no-cache --available \
#     && apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++ curl

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert
COPY --from=builder --chmod=755 /workspace/entry .
COPY --chmod=755 hostname.sh .
COPY modules/backend/envoy-config.yaml .
COPY modules/backend/docker-entry.sh .

ENTRYPOINT ["/docker-entry.sh"]
