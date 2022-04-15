#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export CGO_ENABLED=0
export GO111MODULE=on

.DEFAULT_GOAL := .default

.default: format build lint test

.PHONY: help
help: ## Shows help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.which-go:
	@which go > /dev/null || (echo "install go from https://golang.org/dl/" & exit 1)

.which-ebitenmobile:
	@which ebitenmobile > /dev/null || (echo "install go from https://ebiten.org/documents/mobile.html" & exit 1)

.which-goimports:
	@which goimports > /dev/null || (echo "install goimports from https://pkg.go.dev/golang.org/x/tools/cmd/goimports" & exit 1)

.PHONY: format
format: .which-go .which-goimports ## Formats Go files
	gofmt -s -w $(ROOT)
	goimports -w .

.which-lint:
	@which golangci-lint > /dev/null || (echo "install golangci-lint from https://github.com/golangci/golangci-lint" & exit 1)

.PHONY: lint
lint: .which-lint ## Checks code with Golang CI Lint
	CGO_ENABLED=1 golangci-lint run

.PHONY: build
build: build-darwin build-windows build-wasm

.PHONY: build-darwin
build-darwin: .which-go ## Builds game for macOS
	CGO_ENABLED=1 GOOS=darwin go build -v -o $(ROOT)/shove-it -ldflags="-s -w" $(ROOT)/*.go

.PHONY: build-linux
build-linux: .which-go ## Builds game for linux FIXME: doesn't work
	CGO_ENABLED=1 GOOS=linux go build -v -o $(ROOT)/shove-it -ldflags="-s -w" $(ROOT)/*.go

.PHONY: build-windows
build-windows: .which-go ## Builds game for windows
	GOOS=windows go build -v -o $(ROOT)/shove-it.exe $(ROOT)/*.go

.PHONY: build-wasm
build-wasm: .which-go ## Builds game for Browser
	GOOS=js GOARCH=wasm go build -o shove-it.wasm .
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js .

.PHONY: build-android
build-android: .which-ebitenmobile ## Builds game for Android
	ebitenmobile bind -target android -javapkg net.nasermirzaei89.shove-it -o ./shove-it.aar ./mobile

.PHONY: test
test: .which-go ## Tests go files
	CGO_ENABLED=1 go test -coverpkg=./... -race -coverprofile=./coverage.txt -covermode=atomic $(ROOT)/...
	go tool cover -func coverage.txt
