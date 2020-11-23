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


# Install Terraform
# Versions 0.12.x and 0.13.x are supported
ARG version=0.13.5
RUN curl https://releases.hashicorp.com/terraform/${version}/terraform_${version}_linux_amd64.zip -o /tmp/terraform.zip
RUN unzip /tmp/terraform.zip -d /usr/local/bin/


# Install Go modules
WORKDIR  /srv/terraform-provider-rollbar
COPY go.mod go.sum ./
RUN go mod download -x


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
# Terraform 0.13 `required_providers` syntax is not entirely supported by 0.12,
# so we override it.
RUN ["/bin/bash", "-c", "echo $version; if [[ $version == 0.12* ]]; then mv -v providers012.tf.override providers.tf; fi"]


# Initialize provider
RUN terraform init


# Required environment variable
ENV ROLLBAR_API_KEY=


# Enable trace logging
#ENV TF_LOG=TRACE


# Terraform plan
CMD echo $ROLLBAR_API_KEY && terraform plan -var rollbar_token=$ROLLBAR_API_KEY
