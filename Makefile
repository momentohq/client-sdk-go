.PHONY: install-devtools format imports tidy vet staticcheck lint \
	install-protos-devtools protos build precommit \
	test \
	vendor build-examples run-docs-examples

GOFILES_NOT_NODE = $(shell find . -type f -name '*.go' -not -path "./examples/aws-lambda/infrastructure/*")


install-devtools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.7
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.8.1


format:
	gofmt -s -w ${GOFILES_NOT_NODE}


imports:
	goimports -l -w ${GOFILES_NOT_NODE}


tidy:
	go mod tidy


vet:
	go vet ./...


staticcheck:
	staticcheck ./...


lint: format imports tidy vet staticcheck


install-protos-devtools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


protos:
	@echo "Run `make install-protos-devtools` if you're missing the Go grpc tools."
	@echo "Did you copy momentohq/client_protos/proto/*.proto to internal/protos?"
	# this file is not needed and causes errors, so make sure it's not present
	rm -f internal/protos/httpcache.proto
	protoc -I=internal/protos --go_out=internal/protos --go_opt=paths=source_relative --go-grpc_out=internal/protos --go-grpc_opt=paths=source_relative internal/protos/*.proto


build:
	go build ./...


precommit: lint test


test:
	ginkgo momento/ auth/ batchutils/


vendor:
	go mod vendor


build-examples:
	cd examples && go build ./...


run-docs-examples:
	cd examples && go run docs-examples/main.go
