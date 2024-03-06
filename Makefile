GOFILES_NOT_NODE = $(shell find . -type f -name '*.go' -not -path "./examples/aws-lambda/infrastructure/*")

.PHONY: install-devtools
install-devtools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.8.1

.PHONY: format
format:
	gofmt -s -w ${GOFILES_NOT_NODE}

.PHONY: imports
imports:
	goimports -l -w ${GOFILES_NOT_NODE}

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: lint
lint: format imports tidy vet staticcheck

.PHONY: install-protos-devtools
install-protos-devtools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: protos
protos:
	@echo "Run `make install-protos-devtools` if you're missing the Go grpc tools."
	@echo "Did you copy momentohq/client_protos/protos/*.proto to internal/protos?"
	# this file is not needed and causes errors, so make sure it's not present
	rm -f internal/protos/httpcache.proto
	protoc -I=internal/protos --go_out=internal/protos --go_opt=paths=source_relative --go-grpc_out=internal/protos --go-grpc_opt=paths=source_relative internal/protos/*.proto

.PHONY: build
build:
	go build ./...

.PHONY: precommit
precommit: lint test

.PHONY: test
test:
	ginkgo -v momento/ auth/ batchutils/

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build-examples
build-examples:
	cd examples
	go build ./...
