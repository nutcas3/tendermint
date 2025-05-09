FROM golang:1.24-alpine AS build-env

# Set up dependencies
RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory
WORKDIR /go/src/github.com/serv-chain/serv

# Add source files
COPY . .

# Build the application
RUN go mod download
RUN make build

# Final image
FROM alpine:3.21

# Install necessary packages
RUN apk add --no-cache ca-certificates jq curl bash

WORKDIR /root

# Copy binary from build-env
COPY --from=build-env /go/src/github.com/serv-chain/serv/build/servchaind /usr/bin/servchaind

# Expose ports
# Tendermint P2P, Tendermint RPC, REST API, gRPC
EXPOSE 26656 26657 1317 9090

# Command
CMD ["servchaind", "start"]
