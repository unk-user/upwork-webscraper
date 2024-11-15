# syntax=docker/dockerfile:1

# https://docs.docker.com/language/golang/build-images/
FROM golang:1.23 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
# Remember to build your handler executable for Linux!
# https://github.com/aws/aws-lambda-go/blob/main/README.md#building-your-function
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /main

# Install chromium
FROM public.ecr.aws/lambda/provided:al2 AS chromium

# install brotli, so we can decompress chromium
# we don't have access to brotli out of the box, to install we first need epel
# https://docs.fedoraproject.org/en-US/epel/#what_is_extra_packages_for_enterprise_linux_or_epel
RUN yum -y install amazon-linux-extras && \
    amazon-linux-extras install -y epel && \
    yum -y install brotli wget && \
    wget --progress=dot:giga https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br -O /chromium.br && \
    brotli -d /chromium.br && \
    yum clean all

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2

# install chromium dependencies
RUN yum -y install \
    libX11 \
    nano \
    unzip \
    wget \
    xclock \
    xorg-x11-xauth \
    xterm && \
    yum clean all

# copy in chromium from chromium stage
COPY --from=chromium /chromium /opt/chromium

# grant our program access to chromium
RUN chmod 777 /opt/chromium

# copy in lambda fn from build stage
COPY --from=build /main /main

ENTRYPOINT ["/main"]
CMD ["main.handler"]