FROM golang:1.23-alpine AS build
WORKDIR /app

RUN apk add --no-cache bash

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN env GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-s -w" -o main

FROM alpine:3.18 AS chromium
WORKDIR /tmp

# Install brotli for decompressing Chromium
RUN apk add --no-cache wget brotli && \
    wget --progress=dot:giga https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br -O chromium.br && \
    brotli -d chromium.br && \
    chmod 755 chromium && \
    mv chromium /opt/

FROM public.ecr.aws/lambda/provided:al2023

RUN dnf -y install \
    libX11 \
    nano \
    unzip \
    wget \
    xorg-x11-xauth \
    xterm \
    nss && \
    dnf clean all

COPY --from=chromium /opt/chromium /opt/chromium

RUN chmod 777 /opt/chromium

COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]