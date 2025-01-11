COVERAGE_FILE ?= coverage.out

TARGET ?= gator

.PHONY: build
build:
	@echo "Building ${TARGET}..."
	@mkdir -p .bin
	@go build -o .bin/${TARGET} ./cmd/${TARGET}

.PHONY: test
test:
	@go test -coverpkg='github.com/mashfeii/gator/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	@rm -rf .bin
