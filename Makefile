BINARY  := quine
PKG     := ./cmd/quine
MODULE  := github.com/kehao95/quine

# Build flags
GOFLAGS ?=
LDFLAGS ?=

# Coverage output
COVERPROFILE := coverage.out

.PHONY: all build test test-substrate validate test-v cover cover-html vet clean install audit-model-scenarios check-model-scenarios audit-public-surface check-public-surface audit-control-plane check-control-plane install-githooks

all: vet test-substrate build

build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY) $(PKG)

test:
	@echo "make test is no longer a synonym for go test ./..." >&2
	@echo "Use 'make test-substrate' for deterministic Go tests." >&2
	@echo "Use 'make validate ARGS=\"--change runtime --runtime test_exit_success\"' for layered validation." >&2
	@exit 2

test-substrate:
	go test $(GOFLAGS) -count=1 ./...

validate:
	@test -n "$(ARGS)" || (echo "usage: make validate ARGS='--change runtime --runtime test_exit_success'" >&2; exit 1)
	./tests/validate.sh $(ARGS)

test-v:
	go test $(GOFLAGS) -count=1 -v ./...

cover:
	go test $(GOFLAGS) -count=1 -coverprofile=$(COVERPROFILE) ./...
	go tool cover -func=$(COVERPROFILE)

cover-html: cover
	go tool cover -html=$(COVERPROFILE)

vet:
	go vet ./...

clean:
	rm -f $(BINARY) $(COVERPROFILE)

install:
	go install $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PKG)

audit-model-scenarios:
	./scripts/check-model-scenarios.sh

check-model-scenarios:
	./scripts/check-model-scenarios.sh --strict

audit-public-surface:
	./scripts/check-public-surface.sh

check-public-surface:
	./scripts/check-public-surface.sh

audit-control-plane: audit-model-scenarios audit-public-surface

check-control-plane: check-model-scenarios check-public-surface

install-githooks:
	./scripts/install-githooks.sh
