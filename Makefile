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
	$(GO) build -o bin/net/nettimeout example/net/timeout/timeout.go
	$(GO) build -o bin/pmq/bytes example/pmq/bytes/bytes.go
	$(GO) build -o bin/pmq/mqrequest example/pmq/mqrequest/mqrequest.go
	$(GO) build -o bin/pmq/protobuf example/pmq/protobuf/protobuf.go
	$(GO) build -o bin/pmq/timeout example/pmq/timeout/timeout.go

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

.PHONY: test
test: 
	$(GO) test -v

.PHONY: examples
examples: 
	./bin/net/bytes
	./bin/net/mqrequest
	./bin/net/protobuf
	./bin/net/timeout
	./bin/pmq/bytes
	./bin/pmq/mqrequest
	./bin/pmq/protobuf
	./bin/pmq/timeout

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
