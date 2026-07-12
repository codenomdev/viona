.PHONY: build clean ui

VERSION=development
BIN=codenom
DIR_SRC=./cmd/codenom

GO_ENV=CGO_ENABLED=0 GO111MODULE=on
Revision=$(shell git rev-parse --short HEAD 2>/dev/null || echo "")
GO_FLAGS=-ldflags="-X github.com/codenomdev/viona/cmd.Version=$(VERSION) -X 'github.com/codenomdev/viona/cmd.Revision=$(Revision)' -X 'github.com/codenomdev/viona/cmd.Time=`date +%s`' -extldflags -static"
GO=$(GO_ENV) "$(shell which go)"
TOOLS_BIN := $(shell mkdir -p build/tools && realpath build/tools)

GOLANGCI_VERSION ?= v2.6.2

GOLANGCI = $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)
$(GOLANGCI):
	rm -f $(TOOLS_BIN)/golangci-lint*
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_VERSION)/install.sh | sh -s -- -b $(TOOLS_BIN) $(GOLANGCI_VERSION)
	mv $(TOOLS_BIN)/golangci-lint $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)

# https://dev.to/thewraven/universal-macos-binaries-with-go-1-16-3mm3
universal: generate
	@GOOS=darwin GOARCH=amd64 $(GO_ENV) $(GO) build $(GO_FLAGS) -o ${BIN}_amd64 $(DIR_SRC)
	@GOOS=darwin GOARCH=arm64 $(GO_ENV) $(GO) build $(GO_FLAGS) -o ${BIN}_arm64 $(DIR_SRC)
	@lipo -create -output ${BIN} ${BIN}_amd64 ${BIN}_arm64
	@rm -f ${BIN}_amd64 ${BIN}_arm64

build: generate
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

generate:
	@$(GO) get github.com/google/wire/cmd/wire@v0.5.0
	@$(GO) install github.com/google/wire/cmd/wire@v0.5.0
	@$(GO) generate ./...
	@$(GO) mod tidy

install-ui-packages:
	@corepack enable
	@corepack prepare pnpm@9.7.0 --activate

check:
	@wire flags

ui:
	@cd ui && pnpm pre-install && pnpm build && cd -

lint: generate $(GOLANGCI)

lint-fix: generate $(GOLANGCI)
	@bash ./script/check-asf-header.sh
	$(GOLANGCI) run --fix

# clean all build result
clean:
	@$(GO) clean ./...
	@rm -f $(BIN)

all: clean build
