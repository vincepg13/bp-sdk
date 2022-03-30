PACKAGES=$(shell go list ./... | grep -v '/vendor/')
VERSION := 0.26.0
VERSION := $(VERSION)_$(shell git describe --always)
BUILD_FLAGS = -ldflags "-X github.com/vincepg13/bp-sdk/beyond/vendor/github.com/cosmos/cosmos-sdk/version.Version=${VERSION}"

all: get_tools get_vendor_deps build test

get_tools:
	go get github.com/golang/dep/cmd/dep

build:
	go build $(BUILD_FLAGS) -o bin/beyondcli cmd/beyondcli/main.go && go build $(BUILD_FLAGS) -o bin/beyondd cmd/beyondd/main.go

get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf vendor/
	@dep ensure -v 

test:
	@go test $(PACKAGES)

benchmark:
	@go test -bench=. $(PACKAGES)

.PHONY: all build test benchmark
