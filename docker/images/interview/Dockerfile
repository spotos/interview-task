ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine

RUN set -eux; \
	apk add --no-cache --virtual .build-deps alpine-sdk; \
    go install github.com/githubnemo/CompileDaemon@latest; \
    apk add --no-cache supervisor git gcc musl-dev; \
    apk del .build-deps;

RUN mkdir -p /etc/supervisor/conf.d

WORKDIR /build
