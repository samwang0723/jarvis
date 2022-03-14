# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o /jarvis-api

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates
RUN apk add --no-cache tzdata

WORKDIR /

COPY --from=build_base /app/config.yaml /config.yaml
COPY --from=build_base /jarvis-api /jarvis-api

# This container exposes ports to the outside world
EXPOSE 8080 8081 80 443 22

#USER nonroot:nonroot

# Run the binary program
ENTRYPOINT [ "/jarvis-api" ]
