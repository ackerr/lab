BIN=lab
CHANGE_FILES=`git diff --name-only --diff-filter=AM HEAD | grep --color=never '.go$$' | paste -sd ' ' - || true`

.PHONY: fmt lint test build clean test_all

.DEFAULT: help
help:
	@echo "fmt           : run gofmt"
	@echo "lint          : run golint"
	@echo "test          : run go test"
	@echo "test_all      : run fmt lint test"
	@echo "build         : run go build"
	@echo "clean         : remove the bin"

fmt:
	@gofmt -d -w -e ${CHANGE_FILES}

lint:
	@golint ./...

test:
	@go test ./...

build:
	@go build -o ${BIN}  *.go

clean:
	@git clean -fdx ${BIN}

test_all: fmt lint test build clean
