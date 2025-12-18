# Makefile

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")

build-bin:
	$(info Building awesome binary file...)
	CGO_ENABLED=0 GOOS=linux go build \
 		-ldflags="-X github.com/vendor116/awesome/pkg/version.version=$(VERSION) -s -w" \
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
	oapi-codegen --config api/openapi/models.cfg.yaml api/openapi/openapi.yaml
	@echo "Generating OpenAPI server..."
	oapi-codegen --config api/openapi/server.cfg.yaml api/openapi/openapi.yaml
	@echo "Generating OpenAPI client..."
	oapi-codegen --config api/openapi/client.cfg.yaml api/openapi/openapi.yaml


GOLANGCI_VERSION := v2.7.2

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_VERSION)

OAPI_CODEGEN_VERSION := v2.2.0

install-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

.PHONY: \
	build-bin \
	test \
	lint \
	fix \
	build-docker \
	compose-up \
	gen-openapi \
	install-linter \
	install-oapi-codegen