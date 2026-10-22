GO ?= go
.PHONY: run test vet fmt build demo
run:
	$(GO) run ./cmd/pub
test:
	$(GO) test -race ./...
vet:
	$(GO) vet ./...
fmt:
	$(GO) fmt ./...
build:
	$(GO) build ./...
demo:
	$(GO) run ./cmd/demo
