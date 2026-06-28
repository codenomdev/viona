VERSION=2.0.1
BIN=codenom
DIR_SRC=./cmd/codenom

GO_ENV=CGO_ENABLED=0 GO111MODULE=on
Revision=$(shell git rev-parse --short HEAD 2>/dev/null || echo "")
GO_FLAGS=-ldflags="-X github.com/codenomdev/viona/cmd.Version=$(VERSION) -X 'github.com/codenomdev/viona/cmd.Revision=$(Revision)' -X 'github.com/codenomdev/viona/cmd.Time=`date +%s`' -extldflags -static"
GO=$(GO_ENV) "$(shell which go)"

build: generate
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

generate:
	@$(GO) get github.com/google/wire/cmd/wire@v0.5.0
	@$(GO) install github.com/google/wire/cmd/wire@v0.5.0
	@$(GO) generate ./...
	@$(GO) mod tidy