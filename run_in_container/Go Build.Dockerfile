# First stage: Build the binary
FROM golang:1.20-buster

WORKDIR /app
COPY ./container.yaml ./container.yaml
COPY ./local.yaml ./local.yaml
ENV ENV="container"

EXPOSE 8888