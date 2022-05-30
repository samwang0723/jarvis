.PHONY: test lint

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /third_party/)

lint: lint-check-deps
	@echo "[golangci-lint] linting sources"
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		./...

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint";\
		GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint;\
	fi

migrate:
	@echo "[goose up] do mysql schema migration"
	@goose -dir internal/db/migration mysql "jarvis:password@tcp(0.0.0.0:3306)/jarvis?charset=utf8" up

build:
	@echo "[go build] build executable binary for development"
	@go build -o jarvis-api cmd/main.go

proto:
	@echo "[protoc] generate protobuf related go files, grpc_gateway reversed proxy and swagger"
	@protoc jarvis.v1.proto -I . \
		-I $$GOPATH/src/github.com/samwang0723/jarvis/third_party \
		--go_out ./internal/app/pb --go_opt paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
    	--go-grpc_out ./internal/app/pb --go-grpc_opt paths=source_relative  \
		--grpc-gateway_out ./internal/app/pb/gateway \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt standalone=true \
		--openapiv2_out=logtostderr=true:$$GOPATH/src/github.com/samwang0723/jarvis/api \
		--proto_path=$$GOPATH/src/github.com/samwang0723/jarvis/internal/app/pb

docker-m1:
	@echo "[docker build] build local docker image on Mac M1"
	@docker build -t samwang0723/jarvis-api:m1 -f build/docker/app/Dockerfile.local .

docker-amd64-deps:
	@echo "[docker buildx] install buildx depedency"
	@docker buildx create --name m1-builder
	@docker buildx use m1-builder
	@docker buildx inspect --bootstrap

docker-amd64:
	@echo "[docker buildx] build amd64 version docker image for Ubuntu AWS EC2 instance"
	@docker buildx use m1-builder
	@docker buildx build --load --platform=linux/amd64 -t samwang0723/jarvis-api:latest -f build/docker/app/Dockerfile .
