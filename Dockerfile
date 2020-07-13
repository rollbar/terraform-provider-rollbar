FROM golang:1.13-alpine

# Install neccessary packages.
RUN apk add -U --no-cache bash git
