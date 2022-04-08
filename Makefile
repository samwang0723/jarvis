.PHONY: test lint

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...

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

proto:
	@protoc jarvis.proto -I . \
		--go_out ./pb --go_opt paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
    	--go-grpc_out ./pb --go-grpc_opt paths=source_relative  \
		--grpc-gateway_out ./pb/gateway \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt standalone=true \
		--proto_path=$$GOPATH/src/github.com/samwang0723/jarvis/pb

docker-m1:
	@docker build -t samwang0723/jarvis-api:m1 -f Dockerfile.local .

docker-amd64:
	@docker buildx use m1-builder
	@docker buildx build --load --platform=linux/amd64 -t samwang0723/jarvis-api:latest -f Dockerfile .
