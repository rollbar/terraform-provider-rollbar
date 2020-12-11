FROM golang:alpine
MAINTAINER Jason McVetta <jmcvetta@protonmail.com>


# Update system packages
RUN apk update && apk upgrade --no-cache


# Build folder


# Build dependencies
RUN apk add --no-cache \
        bash \
        curl \
        git \
        make \
        unzip \
        vim


# Enable VI-keys
RUN echo "set -o vi" >> ~/.bashrc


# Install Go modules
WORKDIR  /srv/terraform-provider-rollbar
COPY go.mod go.sum ./
RUN go mod download -x


# Install Terraform
# Versions 0.12.x and 0.13.x are supported
ARG version=0.13.5
RUN curl https://releases.hashicorp.com/terraform/${version}/terraform_${version}_linux_amd64.zip -o /tmp/terraform.zip
RUN unzip /tmp/terraform.zip -d /usr/local/bin/


# Build and install provider
COPY Makefile main.go ./
COPY client client
COPY rollbar rollbar
RUN make build
RUN make install
RUN make install012


# Terraform configuration
RUN mkdir example
COPY example/*.tf example/*.override example/
WORKDIR example/


# Initialize provider
RUN terraform init


# Enable trace logging
#ENV TF_LOG=TRACE


# Terraform plan
ENTRYPOINT ["terraform"]
