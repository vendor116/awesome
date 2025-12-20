# Makefile

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")

build-bin:
	$(info Building awesome binary file...)
	CGO_ENABLED=0 GOOS=linux go build \
 		-ldflags="-X main.version=$(VERSION) -s -w" \
		-o bin/awesome \
		./cmd/awesome

test:
	$(info Running test...)
	go test -v ./...

lint:
	$(info Running golangci-lint lint...)
	golangci-lint run ./...

fix:
	$(info Running golangci-lint fix...)
	golangci-lint run --fix ./...

IMAGE_TAG := vendor116/awesome:latest

build-docker:
	$(info Building Docker image...)
	docker build -f ci/build/Dockerfile -t $(IMAGE_TAG) .

compose-up:
	docker compose -f ./ci/dev/docker-compose.yml up -d

gen-openapi:
	@echo "Generating OpenAPI models..."
	oapi-codegen --config api/openapi/v1/models.cfg.yml api/openapi/v1/openapi.yml
	@echo "Generating OpenAPI server..."
	oapi-codegen --config api/openapi/v1/server.cfg.yml api/openapi/v1/openapi.yml
	@echo "Generating OpenAPI client..."
	oapi-codegen --config api/openapi/v1/client.cfg.yml api/openapi/v1/openapi.yml

gen-protobuf:
	$(info Generating protobuf files...)
	protoc --proto_path=./api/proto \
           --go_out=./pkg/protobuf/awesome \
           --go_opt=paths=source_relative \
           --go-grpc_out=./pkg/protobuf/awesome \
           --go-grpc_opt=paths=source_relative \
           ./api/proto/awesome.proto

GOLANGCI_VERSION := v2.7.2
OAPI_CODEGEN_VERSION := v2.5.1

install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_VERSION)
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: \
	build-bin \
	test \
	lint \
	fix \
	build-docker \
	compose-up \
	gen-openapi \
	gen-protobuf \
	install-tools