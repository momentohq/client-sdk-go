.PHONY: install-devtools install-ginkgo format imports tidy vet staticcheck lint \
	install-protoc-from-client-protos install-protos-devtools update-protos build-protos update-and-build-protos \
	build precommit \
	test test-auth-service test-cache-service test-leaderboard-service test-storage-service test-topics-service \
	vendor build-examples run-docs-examples

GOFILES_NOT_NODE = $(shell find . -type f -name '*.go' -not -path "./examples/aws-lambda/infrastructure/*")
TEST_DIRS = momento/ auth/ batchutils/
GINKGO_OPTS = --no-color -v

install-devtools: install-ginkgo
	@echo "Installing dev tools..."
	go install golang.org/x/tools/cmd/goimports@v0.24.0
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.7


install-ginkgo:
	@echo "Installing ginkgo..."
	@go install github.com/onsi/ginkgo/v2/ginkgo@v2.8.1

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

install-protoc-from-client-protos:
	@echo "Installing protoc from latest client_protos release..."
	@temp_dir=$$(mktemp -d) && \
		latest_tag=$$(git ls-remote --tags --sort="v:refname" https://github.com/momentohq/client_protos.git | tail -n 1 | sed 's!.*/!!') && \
		echo "Latest release tag: $$latest_tag" && \
		git -c advice.detachedHead=false clone --branch "$$latest_tag" https://github.com/momentohq/client_protos.git $$temp_dir && \
		cd $$temp_dir && \
		./install_protoc.sh && \
		rm -rf $$temp_dir

install-protos-devtools:
	@echo "Installing protos dev tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0

update-protos:
	@echo "Updating .proto files from the latest release of the client_protos repository..."
# Note: httpcache.proto is not needed and causes errors, so make sure it's not present
	@temp_dir=$$(mktemp -d) && \
		latest_tag=$$(git ls-remote --tags --sort="v:refname" https://github.com/momentohq/client_protos.git | tail -n 1 | sed 's!.*/!!') && \
		echo "Latest release tag: $$latest_tag" && \
		git -c advice.detachedHead=false clone --branch "$$latest_tag" https://github.com/momentohq/client_protos.git $$temp_dir && \
		cp $$temp_dir/proto/*.proto internal/protos/ && \
		rm -f internal/protos/httpcache.proto && \
		rm -rf $$temp_dir

build-protos:
	@echo "Generating go code from protos..."
	@echo "Run `make install-protos-devtools` if you're missing the Go grpc tools."
	protoc -I=internal/protos --go_out=internal/protos --go_opt=paths=source_relative --go-grpc_out=internal/protos --go-grpc_opt=paths=source_relative internal/protos/*.proto

update-and-build-protos: update-protos build-protos

build:
	@echo "Building..."
	go build ./...


precommit: lint test


test: install-ginkgo
	@echo "Running tests..."
	@ginkgo ${GINKGO_OPTS} ${TEST_DIRS}

prod-test: install-ginkgo
	@echo "Running tests with consistent reads..."
	@CONSISTENT_READS=1 ginkgo ${GINKGO_OPTS} ${TEST_DIRS}


test-auth-service: install-ginkgo
	@echo "Testing auth service..."
	@CONSISTENT_READS=1 ginkgo ${GINKGO_OPTS} --label-filter auth-service ${TEST_DIRS}


test-cache-service: install-ginkgo
	@echo "Testing cache service..."
	@CONSISTENT_READS=1 ginkgo ${GINKGO_OPTS} --label-filter cache-service ${TEST_DIRS}


test-leaderboard-service: install-ginkgo
	@echo "Testing leaderboard service..."
	@CONSISTENT_READS=1 ginkgo ${GINKGO_OPTS} --label-filter leaderboard-service ${TEST_DIRS}


test-storage-service: install-ginkgo
	@echo "Testing storage service..."
	@CONSISTENT_READS=1 ginkgo ${GINKGO_OPTS} --label-filter storage-service ${TEST_DIRS}


test-topics-service: install-ginkgo
	@echo "Testing topics service..."
	@CONSISTENT_READS=1 ginkgo ${GINKGO_OPTS} --label-filter topics-service ${TEST_DIRS}


vendor:
	@echo "Vendoring..."
	go mod vendor


build-examples:
	@echo "Building examples..."
	cd examples && go build ./...


run-docs-examples:
	@echo "Running docs examples..."
	cd examples && go run docs-examples/main.go
