# syntax=docker/dockerfile:1
FROM golang:1.17-alpine AS build_base

# Add Maintainer Info
LABEL maintainer="Sam Wang <sam.wang.0723@gmail.com>"

RUN apk add --no-cache git
RUN apk update && apk add ca-certificates && apk add tzdata

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Unit tests
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /jarvis-api

# Start fresh from a smaller image
FROM scratch

#RUN apk add ca-certificates
#RUN apk add --no-cache tzdata

WORKDIR /

COPY --from=build_base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build_base /app/config.prod.yaml /config.prod.yaml
COPY --from=build_base /jarvis-api /jarvis-api

# This container exposes ports to the outside world
EXPOSE 8080 8081 80 443 22 3306 6379

#USER nonroot:nonroot
#RUN chmod +x /jarvis-api
ENV TZ=Asia/Taipei

# Run the binary program
ENTRYPOINT [ "/jarvis-api" ]

#docker buildx build --load --platform=linux/amd64 -t samwang0723/jarvis-api:latest -f Dockerfile .
