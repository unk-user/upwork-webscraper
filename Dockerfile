FROM golang:1.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN env GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o main

FROM amazonlinux:2023 AS chromium

RUN yum -y install wget brotli && \
    wget --progress=dot:giga https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br -O /chromium.br && \
    brotli -d /chromium.br && \
    yum clean all

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

COPY --from=chromium /chromium /opt/chromium

RUN chmod 777 /opt/chromium

COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]