FROM golang:1.10.2-alpine

# Install neccessary packages.
RUN apk add -U --no-cache bash git
