BINARY  := quine
PKG     := ./cmd/quine
MODULE  := github.com/kehao95/quine

# Panopticon
PAN_BINARY := panopticon
PAN_PKG    := ./cmd/panopticon

# Build flags
GOFLAGS ?=
LDFLAGS ?=

# Coverage output
COVERPROFILE := coverage.out

.PHONY: all build test test-v cover cover-html vet clean install panopticon panopticon-web panopticon-dev

all: vet test build

build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY) $(PKG)

panopticon:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(PAN_BINARY) $(PAN_PKG)

panopticon-web:
	cd web && npm run build

panopticon-dev:
	@echo "Starting Panopticon backend on :8900 and frontend on :3000..."
	@go run $(PAN_PKG) &
	@cd web && npm run dev

test:
	go test $(GOFLAGS) -count=1 ./...

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
	rm -f $(BINARY) $(PAN_BINARY) $(COVERPROFILE)

install:
	go install $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PKG)
