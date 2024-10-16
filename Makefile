#!/usr/bin/make -f

# for dockerized protobuf tools
DOCKER := $(shell which docker)

export GO111MODULE = on

all: test

install:
	$(MAKE) -C demo install

build:
	$(MAKE) -C demo build

build-linux-static:
	$(MAKE) -C demo build-linux-static

build-docker:
	$(DOCKER) build --tag babylonlabs-io/bcd -f Dockerfile \
		$(shell git rev-parse --show-toplevel)

########################################
### Testing

test-all: test

test:
	$(MAKE) -C demo test
	$(MAKE) -C x test
	$(MAKE) -C tests/e2e/ test

test-e2e:
	$(MAKE) -C tests/e2e/ test

###############################################################################
###                                Linting                                  ###
###############################################################################

format-tools:
	go install mvdan.cc/gofumpt@v0.4.0
	go install github.com/client9/misspell/cmd/misspell@v0.3.4
	go install github.com/daixiang0/gci@v0.11.2

lint: format-tools
	$(MAKE) -C demo lint
	$(MAKE) -C tests/e2e lint
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./x/vendor*" -not -path "*.git*" -not -path "*_test.go" | xargs gofumpt -d -s

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./x/vendor*" -not -path "./contracts*" -not -path "./packages*" -not -path "./docs*"| xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./x/vendor*" -not -path "./contracts*" -not -path "./packages*" -not -path "./docs*"| xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./tests/system/vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gci write --skip-generated -s standard -s default -s "prefix(cosmossdk.io)" -s "prefix(github.com/cosmos/cosmos-sdk)" -s "prefix(github.com/babylonlabs-io/babylon-sdk)" --custom-order


###############################################################################
###                                Protobuf                                 ###
###############################################################################
protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-swagger-gen:
	@./scripts/protoc-swagger-gen.sh

proto-lint:
	@$(protoImage) buf lint --error-format=json

.PHONY: all install \
	build build-linux-static test test-all test-e2e \
	proto-all proto-format proto-swagger-gen proto-lint
