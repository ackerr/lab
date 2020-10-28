VERSION=0.2.5

BIN=lab
DIST_DIR=dist
TEMP_DIR=temp

.PHONY: fmt lint test build clean test_all

.DEFAULT: help
help:
	@echo "fmt           : run gofmt"
	@echo "lint          : run golangci-lint"
	@echo "build         : run go build"
	@echo "clean         : remove the bin"
	@echo "install       : install the dependence"
	@echo "test          : run go test"
	@echo "coverage      : run go test, and collection coverage"
	@echo "coverage_html : run go test, and show the coverage"
	@echo "test_all      : run lint test"

fmt:
	@gofmt -d -w -e .

lint:
	@go vet ./...
	@go mod tidy
	@golangci-lint run

install:
	@go mod download

build:
	@go build -v -o ${BIN} *.go

clean:
	@git clean -fdx ${BIN} ${DIST_DIR}

.PHONY: test coverage

COVERAGE_FILE=".coverage"
GITLAB_TOKEN="gitlab token"
GITLAB_BASE_URL="gitlab base url"

test:
	@export $GITLAB_TOKEN
	@export $GITLAB_BASE_URL
	@mkdir -p ${TEMP_DIR}
	@ROOT=${PWD} go test -cover -covermode=atomic -coverpkg=./... -coverprofile=${COVERAGE_FILE} -v -race ./...
	@rm -rf ${TEMP_DIR}

coverage:
	go tool cover -func ${COVERAGE_FILE}

coverage_html:
	go tool cover -html ${COVERAGE_FILE}

test_all: lint test

release:
	@bumpversion patch
