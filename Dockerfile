FROM golang:1.10.2-alpine

# Install the neccessary packages for bulding the binary.
RUN apk add -U --no-cache bash git
