.PHONY: install dev-server dev-web format format-check lint test test-go test-web audit build check

GO ?= go
GOFMT ?= gofmt
NPM ?= npm

install:
	$(NPM) --prefix frontend ci

dev-server:
	$(GO) run ./cmd/server

dev-web:
	$(NPM) --prefix frontend run dev

format:
	$(GOFMT) -w cmd internal
	$(NPM) --prefix frontend exec -- eslint . --fix

format-check:
	@test -z "$$($(GOFMT) -l cmd internal)" || ($(GOFMT) -l cmd internal && exit 1)

lint: format-check
	$(GO) vet ./cmd/... ./internal/...
	$(NPM) --prefix frontend run lint
	$(NPM) --prefix frontend run typecheck

test: test-go test-web

test-go:
	$(GO) test -race ./cmd/... ./internal/...

test-web:
	$(NPM) --prefix frontend test

audit:
	$(NPM) --prefix frontend audit

build:
	$(NPM) --prefix frontend run build
	mkdir -p bin
	$(GO) build -o bin/gotopia ./cmd/server

check: lint test build
