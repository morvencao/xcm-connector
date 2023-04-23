.DEFAULT_GOAL :=help

SHELL :=/bin/bash

VERSION ?=$(shell date +%s)
IMAGE_TAG ?=$(VERSION)
IMAGE_REGISTRY ?=quay.io/morvencao
IMAGE ?=xcm-connector
IMAGE_NAME ?=$(IMAGE_REGISTRY)/$(IMAGE):$(IMAGE_TAG)

KUBECTL ?=kubectl
GOLANGCI_LINT_BIN :=$(shell go env GOPATH)/bin/golangci-lint

GIT_HOST ?=github.com/morvencao/xcm-connector
BASE_DIR :=$(shell basename $(PWD))
DEST :=$(GOPATH)/src/$(GIT_HOST)/$(BASE_DIR)

# Add packages to do unit test
GO_TEST_PACKAGES :=./pkg/...

# Prints a list of useful targets.
help:
	@echo ""
	@echo "xCM Connector"
	@echo ""
	@echo "make verify               verify source code"
	@echo "make lint                 run golangci-lint"
	@echo "make build                build binaries"
	@echo "make test                 run unit tests"
	@echo "make image                build docker image"
	@echo "make push                 push docker image"
	@echo "make deploy               deploy via templates to local openshift instance"
	@echo "make undeploy             undeploy from local openshift instance"
	@echo "make clean                delete temporary generated files"
	@echo "$(fake)"
.PHONY: help

# Verifies that source passes standard checks.
verify:
	go vet \
		./cmd/... \
		./pkg/...
	! gofmt -l cmd pkg |\
		sed 's/^/Unformatted file: /' |\
		grep .
.PHONY: verify

# Runs our linter to verify that everything is following best practices
# Requires golangci-lint to be installed @ $(go env GOPATH)/bin/golangci-lint
lint:
	$(GOLANGCI_LINT_BIN) run \
		./cmd/... \
		./pkg/...
.PHONY: lint

# Build binaries
build:
	go build ./cmd/xcm-connector
.PHONY: build

# Runs the unit tests.
#
# Args:
#   TESTFLAGS: Flags to pass to `go test`. The `-v` argument is always passed.
#
# Examples:
#   make test TESTFLAGS="-run TestSomething"
test:
	go test --format short-verbose -- -p 1 -v $(TESTFLAGS) \
		./pkg/... \
		./cmd/...
.PHONY: test

image:
	docker build -t "$(IMAGE_NAME)" .
.PHONY: image

push: image
	docker push "$(IMAGE_NAME)"
.PHONY: push

deploy:
	kubectl apply -k ./deploy
.PHONY: deploy

undeploy:
	kubectl delete -k ./deploy
.PHONY: undeploy

# Delete temporary files
clean:
	rm -rf \
		xcm-connector \
.PHONY: clean
