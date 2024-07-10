# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build_base

# Add Maintainer Info
LABEL maintainer="Sam Wang <sam.wang.0723@gmail.com>"

# Accept build-time argument for environment
ARG ENVIRONMENT=dev

# Set environment variable
ENV ENVIRONMENT=${ENVIRONMENT}

RUN apk add --no-cache git

# Conditionally execute commands based on the environment
RUN if [ "$ENVIRONMENT" = "prod" ]; then \
  echo "Installing production-specific certification and tzdata"; \
  RUN apk update && apk add ca-certificates && apk add tzdata \
fi

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v -tags all_tests -coverprofile=coverage.txt -covermode=atomic $(go list ./... | grep -v /third_party/)

# Build the Go app
RUN go build -o /jarvis-api cmd/main.go

# Start fresh from a smaller image
FROM scratch

WORKDIR /

# Set environment variable
ENV ENVIRONMENT=${ENVIRONMENT}

# Conditionally execute commands based on the environment
RUN if [ "$ENVIRONMENT" = "prod" ]; then \
  echo "Copy production-specific certification and tzdata"; \
  COPY --from=build_base /usr/share/zoneinfo /usr/share/zoneinfo \
  COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ \
fi

COPY --from=build_base /app/configs/config.${ENVIRONMENT}.yaml /configs/config.${ENVIRONMENT}.yaml
COPY --from=build_base /jarvis-api /jarvis-api

# This container exposes ports to the outside world
EXPOSE 8080 8081 80 443 22 3306 6379

#USER nonroot:nonroot
#RUN chmod +x /jarvis-api
ENV TZ=Asia/Taipei

# Run the binary program
ENTRYPOINT [ "/jarvis-api" ]

###############
# development #
###############

# docker build -t samwang0723/jarvis-api:m1 --build-arg ENVIRONMENT=local -f Dockerfile.local .
# docker run -p 8080-8081:8080-8081 --name jarvis-api-container --env-file=./.env samwang0723/jarvis-api:m1

###############
# production  #
###############

# docker buildx build --load --platform=linux/amd64 -t samwang0723/jarvis-api:latest --build-arg ENVIRONMENT=prod -f Dockerfile .
