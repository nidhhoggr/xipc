GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

all: build

.PHONY: build
build: 
	$(GO) mod tidy
	$(GO) build -o bin/net/bytes example/net/bytes/bytes.go
	$(GO) build -o bin/net/mqrequest example/net/mqrequest/mqrequest.go
	$(GO) build -o bin/net/protobuf example/net/protobuf/protobuf.go
	$(GO) build -o bin/net/timeout example/net/timeout/timeout.go
	$(GO) build -o bin/pmq/bytes example/pmq/bytes/bytes.go
	$(GO) build -o bin/pmq/mqrequest example/pmq/mqrequest/mqrequest.go
	$(GO) build -o bin/pmq/protobuf example/pmq/protobuf/protobuf.go
	$(GO) build -o bin/pmq/timeout example/pmq/timeout/timeout.go
	$(GO) build -o bin/mem/bytes example/mem/bytes/bytes.go
	$(GO) build -o bin/mem/mqrequest example/mem/mqrequest/mqrequest.go
	$(GO) build -o bin/mem/protobuf example/mem/protobuf/protobuf.go
	$(GO) build -o bin/mem/timeout example/mem/timeout/timeout.go


.PHONY: run
run: 
	$(GO) mod tidy
	$(GO) run --race example/net/bytes/bytes.go
	$(GO) run --race example/net/mqrequest/mqrequest.go
	$(GO) run --race example/net/protobuf/protobuf.go
	$(GO) run --race example/net/timeout/timeout.go
	$(GO) run --race example/pmq/bytes/bytes.go
	$(GO) run --race example/pmq/mqrequest/mqrequest.go
	$(GO) run --race example/pmq/protobuf/protobuf.go
	$(GO) run --race example/pmq/timeout/timeout.go
	$(GO) run --race example/mem/bytes/bytes.go
	$(GO) run --race example/mem/mqrequest/mqrequest.go
	$(GO) run --race example/mem/protobuf/protobuf.go
	$(GO) run --race example/mem/timeout/timeout.go

.PHONY: test
test: 
	$(GO) test -v

protogen:
	protoc protos/*.proto  --go_out=.
	protoc example/protos/*.proto  --go_out=.

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
  if [ -n "$$diff" ]; then \
    echo "Please run 'make fmt' and commit the result:"; \
    echo "$${diff}"; \
    exit 1; \
  fi;
