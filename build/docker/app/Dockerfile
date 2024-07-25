ARG GO_VERSION=1.22.5

###########
# MODULES #
###########

FROM --platform=${BUILDPLATFORM} public.ecr.aws/docker/library/golang:${GO_VERSION}-alpine AS modules

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN apk add --no-cache git ca-certificates tzdata

RUN go mod download

###########
# BUILDER #
###########

FROM --platform=${BUILDPLATFORM} public.ecr.aws/docker/library/golang:${GO_VERSION}-alpine AS builder

COPY --from=modules /go/pkg /go/pkg

RUN adduser -D -u 10001 nonroot

WORKDIR /src

COPY ./ ./

# Add Maintainer Info
LABEL maintainer="Sam Wang <sam.wang.0723@gmail.com>"

ARG COMPILATION_MAIN_FILES
ARG GLOBAL_VAR_PKG

ARG TARGETOS
ARG TARGETARCH

ARG LAST_MAIN_COMMIT_HASH
ARG LAST_MAIN_COMMIT_TIME
ARG RELEASE_TAG

ENV FLAG="-X ${GLOBAL_VAR_PKG}.CommitTime=${LAST_MAIN_COMMIT_TIME}"
ENV FLAG="$FLAG -X ${GLOBAL_VAR_PKG}.CommitHash=${LAST_MAIN_COMMIT_HASH}"
ENV FLAG="$FLAG -X ${GLOBAL_VAR_PKG}.ReleaseTag=${RELEASE_TAG}"

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -installsuffix 'static' \
    -ldflags "-s -w $FLAG" \
    -buildvcs=true \
    -o /app ${COMPILATION_MAIN_FILES}

#########
# FINAL #
#########

FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=modules /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=modules /usr/share/zoneinfo /usr/share/zoneinfo

COPY ./configs /configs

COPY --from=builder /app /app

ENV TZ=Asia/Taipei

USER nonroot

CMD ["/app"]
