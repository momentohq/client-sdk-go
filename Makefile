.PHONY: install-devtools format imports tidy vet staticcheck lint \
	install-protos-devtools protos build precommit \
	test \
	vendor build-examples run-docs-examples

GOFILES_NOT_NODE = $(shell find . -type f -name '*.go' -not -path "./examples/aws-lambda/infrastructure/*")


install-devtools:
	@echo "Installing dev tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.7
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.8.1


format:
	@echo "Formatting code..."
	gofmt -s -w ${GOFILES_NOT_NODE}


imports:
	@echo "Running goimports..."
	goimports -l -w ${GOFILES_NOT_NODE}


tidy:
	@echo "Running go mod tidy..."
	go mod tidy


vet:
	@echo "Running go vet..."
	go vet ./...


staticcheck:
	@echo "Running staticcheck..."
	staticcheck ./...


lint: format imports tidy vet staticcheck


install-protos-devtools:
	@echo "Installing protos dev tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


protos:
	@echo "Generating protos..."
	@echo "Run `make install-protos-devtools` if you're missing the Go grpc tools."
	@echo "Did you copy momentohq/client_protos/proto/*.proto to internal/protos?"
	# this file is not needed and causes errors, so make sure it's not present
	rm -f internal/protos/httpcache.proto
	protoc -I=internal/protos --go_out=internal/protos --go_opt=paths=source_relative --go-grpc_out=internal/protos --go-grpc_opt=paths=source_relative internal/protos/*.proto


build:
	@echo "Building..."
	go build ./...


precommit: lint test


test:
	@echo "Running tests..."
	ginkgo momento/ auth/ batchutils/


vendor:
	@echo "Vendoring..."
	go mod vendor


build-examples:
	@echo "Building examples..."
	cd examples && go build ./...


run-docs-examples:
	@echo "Running docs examples..."
	cd examples && go run docs-examples/main.go
