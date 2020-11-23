FROM golang:alpine
MAINTAINER Jason McVetta <jmcvetta@protonmail.com>


# Update system packages
RUN apk update && apk upgrade --no-cache


# Terraform version
ARG version=0.12.29


# Build folder
ENV buildfolder=/srv/terraform-provider-rollbar


# Build dependencies
RUN apk add --no-cache \
        bash \
        curl \
        git \
        make \
        unzip \
        vim

RUN echo "set -o vi" >> ~/.bashrc


# Install Terraform
RUN curl https://releases.hashicorp.com/terraform/${version}/terraform_${version}_linux_amd64.zip -o /tmp/terraform.zip
RUN unzip /tmp/terraform.zip -d /usr/local/bin/


# Install Go modules
WORKDIR ${buildfolder}
COPY go.mod go.sum ./
RUN go mod download -x


# Build provider
COPY Makefile main.go ./
COPY client client
COPY rollbar rollbar
RUN make build
ENV TF_LOG=TRACE
RUN make install012


# Test provider
RUN mkdir example
COPY example/main.tf example/
COPY example/providers012.tf.example example/providers012.tf
RUN make plan
