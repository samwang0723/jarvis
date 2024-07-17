.PHONY: test lint bench lint-skip-fix migrate proto build build-docker install vendor deploy rollback

help: ## show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

PROJECT_NAME?=jarvis
APP_NAME_UND=$(shell echo "$(PROJECT_NAME)" | tr '-' '_')

APP_NAME?=jarvis-api
VERSION?=v2.0.1
SQLC_VERSION=v1.26.0
GCI_VERSION=0.11.2
TRIVY_VERSION=0.37.1
GOLANG_LINTER_VERSION=1.59.1
GO_VERSION=$(shell cat .go_version)

SHELL = /bin/bash
SOURCE_LIST = $$(go list ./... | grep -v /third_party/ | grep -v /internal/app/pb | grep -v /cmd | grep -v /internal/cache/mocks | grep -v /internal/db/main/sqlc | grep -v /database | grep -v /internal/cronjob/mocks | grep -v /internal/services/mocks | grep -v /internal/kafka/mocks )

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

tool-version-check:
	@( \
	    INSTALLED_TOOL_VERSION=$$($(tool_version_check) | grep $(tool_version)); \
		if [ -z "$$INSTALLED_TOOL_VERSION" ]; then \
			echo "$(tool_version_check) mismatch $(tool_version)"; \
			echo "INSTALLED_TOOL_VERSION: $$INSTALLED_TOOL_VERSION"; \
			echo "current version: $$($(tool_version_check))"; \
			exit 1; \
		fi \
	)

###########
# install #
###########
## install: Install go dependencies
install:
	go mod tidy
	go mod download
	go get ./...

# vendor: Vendor go modules
vendor:
	go mod vendor

########
# test #
########

test: test-race test-leak test-coverage-ci ## launch all tests

test-race: ## launch all tests with race detection
	go test $(SOURCE_LIST)  -cover -race -count=1

test-leak: ## launch all tests with leak detection (if possible)
	go test $(SOURCE_LIST)  -leak

test-coverage-ci:
	go test -v $(SOURCE_LIST) -cover -race -covermode=atomic -coverprofile=coverage.out

test-coverage-report:
	go test -v $(SOURCE_LIST) -cover -race -covermode=atomic -coverprofile=coverage.out
	go tool cover -html=coverage.out

########
# lint #
########

lint: lint-check-deps ## lints the entire codebase
	@make tool-version-check tool_version_check="golangci-lint version" tool_version=$(GOLANG_LINTER_VERSION)
	@golangci-lint run ./... --config=./.golangci.yaml --timeout=15m && \
	if [ $$(gofumpt -e -l --extra cmd/ | wc -l) = "0" ] && \
		[ $$(gofumpt -e -l --extra internal/ | wc -l) = "0" ] && \
		[ $$(gofumpt -e -l --extra configs/ | wc -l) = "0" ] ; \
		then exit 0; \
	else \
		echo "these files needs to be gofumpt-ed"; \
		gofumpt -e -l --extra cmd/; \
		gofumpt -e -l --extra internal/; \
		gofumpt -e -l --extra configs/; \
	fi

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint";\
		GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint;\
	fi

lint-skip-fix: ## skip linting the system generate files
	@git checkout head internal/app/pb
	@git checkout head third_party/

###########
#   GCI   #
###########

gci-format: ## format repo through gci linter
	@make tool-version-check tool_version_check="gci --version" tool_version=${GCI_VERSION}
	gci write --skip-generated -s standard -s default -s "Prefix(github.com/samwang0723)" -s "Prefix(jarvis)" ./

#############
# benchmark #
#############

bench: ## launch benches
	go test $(SOURCE_LIST) -bench=. -benchmem | tee ./bench.txt

bench-compare: ## compare benches results
	benchstat ./bench.txt

#######
# sec #
#######

sec-scan: trivy-scan vuln-scan ## scan for security and vulnerabilities

trivy-scan: ## scan for sec issues with trivy (trivy binary needed)
	@make tool-version-check tool_version_check="trivy --version" tool_version=$(TRIVY_VERSION)
	trivy fs --exit-code 1 --no-progress --severity CRITICAL ./

vuln-scan: ## scan for vuln issues with trivy (trivy binary needed)
	govulncheck ./...

###########
#  mock   #
###########

mock-gen: ## generate mocks
	go generate $(SOURCE_LIST)

############
#   sqlc   #
############

sqlc: ## gen sqlc code for your app
	@make tool-version-check tool_version_check="sqlc version" tool_version=$(SQLC_VERSION)
	sqlc generate -f ./database/sqlc/sqlc.yaml

###########
# migrate #
###########

migrate: ## migrate db to latest version
	go run ./cmd/migrate

db-pg-init-main: ## create users and passwords in postgres for your app
	@( \
	printf "Enter host for db(localhost): \n"; read -rs DB_HOST &&\
	printf "Enter pass for db: \n"; read -rs DB_PASSWORD &&\
	printf "Enter port(5432...): \n"; read -r DB_PORT &&\
	sed \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/APP_NAME_UND/$(APP_NAME_UND)/g" \
	./database/init/init.sql | \
	PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U postgres -f - \
	)

db-pg-migrate:
	@( \
	printf "Enter host for db(localhost): \n"; read -rs DB_HOST &&\
	printf "Enter pass for db: \n"; read -rs DB_PASSWORD &&\
	printf "Enter port(5432...): \n"; read -r DB_PORT &&\
	sed \
	-e "s/DB_HOST/$$DB_HOST/g" \
	-e "s/DB_PORT/$$DB_PORT/g" \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/APP_NAME_UND/$(APP_NAME_UND)/g" \
	./database/migrations/main.go > ./database/migrations/tmp.go && \
	go run ./database/migrations/tmp.go up && \
	rm ./database/migrations/tmp.go \
	)


#############
#  upgrade  #
############

upgrade: ## upgrade all dependencies, dangerous!!
	go mod tidy && \
	go get -t -u ./... && \
	go mod tidy

NEW_VERSION = "default"
upgrade-go: ## upgrade go version(example: make upgrade-go NEW_VERSION="1.22.5")
	sed -i '' "s/$(GO_VERSION)/$(NEW_VERSION)/" .go_version && \
	go mod edit -go $(shell echo $(NEW_VERSION) | cut -d. -f1,2) && \
	go mod tidy

#########
# PROTO #
#########

# protoc -I third_party --openapiv2_out api --openapiv2_opt logtostderr=true --proto_path=internal/app/pb jarvis.v1.proto
proto: ## generate proto files
	@echo "[protoc] generate protobuf related go files, grpc_gateway reversed proxy"
	@protoc jarvis.v1.proto -I . \
		-I ./third_party \
		-I ./internal/app/pb \
		--go_out ./internal/app/pb --go_opt paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
    	--go-grpc_out ./internal/app/pb --go-grpc_opt paths=source_relative  \
		--grpc-gateway_out ./internal/app/pb/gateway \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt standalone=true

	@echo "[protoc] generate openapiv2 swagger json"
	@protoc -I ./third_party -I ./internal/app/pb --openapiv2_out api --openapiv2_opt logtostderr=true jarvis.v1.proto



##############
#   build    #
##############
build: lint test bench sec-scan docker-build ## lint, test, bench and sec scan before building the docker image
	@printf "\nyou can now deploy to your env of choice:\ncd deploy\nENV=dev make deploy-latest\n"

LAST_MAIN_COMMIT_HASH=$(shell git rev-parse --short HEAD)
LAST_MAIN_COMMIT_TIME=$(shell git show --no-patch --format=%cd --date=iso-strict HEAD)
RELEASE_TAG=$(shell git describe --abbrev=0 --tags)
REPO_NAME=samwang0723

docker-build: lint test docker-build-api ## build docker image in M1 device
	@printf "\nyou can now deploy to your env of choice:\ncd deploy\nENV=dev make deploy-latest\n"

docker-build-api: TAG_NAME=$(REPO_NAME)/$(APP_NAME) ## docker build for api
docker-build-api: COMPILATION_MAIN_FILES="./cmd/api/*.go"
docker-build-api: GLOBAL_VAR_PKG="main"
docker-build-api:
	$(call docker-build-generic)

# use a function instead of a target for docker-build-generic
# since it's impossible for a target to be run twice in a make call
define docker-build-generic
	if [ -n "${PLATFORM}" ]; then \
		PLATFORM_FLAG="--platform ${PLATFORM}"; \
	fi; \
	docker run --privileged --rm tonistiigi/binfmt --install all; \
	DOCKER_BUILDKIT=1 \
	docker buildx build \
		-f build/docker/app/Dockerfile \
		-t $(TAG_NAME) \
		$$PLATFORM_FLAG \
		--build-arg COMPILATION_MAIN_FILES=$(COMPILATION_MAIN_FILES) \
		--build-arg GLOBAL_VAR_PKG=$(GLOBAL_VAR_PKG) \
		--build-arg LAST_MAIN_COMMIT_HASH=$(LAST_MAIN_COMMIT_HASH) \
		--build-arg LAST_MAIN_COMMIT_TIME=$(LAST_MAIN_COMMIT_TIME) \
		--build-arg RELEASE_TAG=$(RELEASE_TAG) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--ssh default \
		--progress=plain \
		--load \
		./
endef

##################
# k8s Deployment #
##################
deploy:
	@kubectl apply -f deployments/helm/jarvis/deployment.yaml
	@kubectl rollout status deployment/jarvis-api

rollback:
	@kubectl rollout undo deployment/jarvis-api

#############
# changelog #
#############

MOD_VERSION = $(shell git describe --abbrev=0 --tags `git rev-list --tags --max-count=1`)

MESSAGE_CHANGELOG_COMMIT="chore(changelog): update CHANGELOG.md for $(MOD_VERSION)"

changelog-gen: ## generates the changelog in CHANGELOG.md
	@git cliff -o ./CHANGELOG.md && \
	printf "\nchangelog generated!\n"
	git add CHANGELOG.md

changelog-commit:
	git commit -m $(MESSAGE_CHANGELOG_COMMIT) ./CHANGELOG.md
