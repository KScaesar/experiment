# First stage: Build the binary
FROM golang:1.20-buster AS builder

LABEL stage=builder
WORKDIR /app

COPY . .
RUN go mod download && CGO_ENABLED=0 go build -trimpath -o ./server ./...

# Second stage: Copy the binary from the builder stage and run it
FROM alpine:latest

WORKDIR /app
ENV ENV="container"

COPY --from=builder /app/server .
COPY ./container.yaml ./container.yaml

EXPOSE 8888
CMD ["./server"]