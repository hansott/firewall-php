# Docker container used for building Zen for PHP from source on Ubuntu

FROM --platform=linux/amd64 ubuntu:20.04

ARG PHP_VERSION=8.1


RUN apt-get update
RUN apt install software-properties-common -y
RUN add-apt-repository ppa:ondrej/php -y
RUN apt update
RUN apt install php8.1 php8.1-cli php8.1-fpm php8.1-dev php8.1-xml php8.1-curl php8.1-mbstring php8.1-zip php8.1-gd php8.1-mysql php8.1-soap php8.1-intl -y
RUN apt-get install -y wget autoconf bison re2c libxml2-dev libssl-dev libcurl4-gnutls-dev protobuf-compiler protobuf-compiler-grpc git
RUN wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="${HOME}/go"
ENV PATH="${PATH}:${GOPATH}/bin"
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
