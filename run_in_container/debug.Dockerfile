# First stage: Build the binary
FROM golang:1.20-buster AS builder

LABEL stage=builder
WORKDIR /app

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .
RUN go mod download && CGO_ENABLED=0 go build -gcflags="all=-N -l" -trimpath -o ./server ./...

# Second stage: Copy the binary from the builder stage and run it
FROM golang:1.20-buster

WORKDIR /app
ENV ENV="container"

COPY --from=builder /go/bin/dlv .
COPY --from=builder /app/server .
COPY ./container.yaml ./container.yaml

EXPOSE 8888 40000
#CMD ["./server"]
#CMD ["ls","-alh"]
CMD ["./dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./server"]