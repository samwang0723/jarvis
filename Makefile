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
