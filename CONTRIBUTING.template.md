{{ ossHeader }}

# Welcome to client-sdk-go contributing guide :wave:

Thank you for taking your time to contribute to our Go SDK!
This guide will provide you information to start your own development and testing.
Happy coding :dancer:

## Submitting

If you've found a bug, or have a suggestion, please [open an issue in our project](https://github.com/momentohq/client-sdk-go/issues).

If you want to submit a change, please [submit a pull request to our project](https://github.com/momentohq/client-sdk-go/pulls). Use the normal [Github pull request process](https://docs.github.com/en/pull-requests). Please run `make precommit` before submitting your pull request; see below for more information.

## Minimum Go version

Our minimum supported Go version is currently `1.19`. You can download it from [go.dev](https://go.dev/).

## Requirements

To make development easier, we provide a [Makefile](https://golangdocs.com/makefiles-golang) to do common development tasks. If you're on Windows, you can get `make` by installing [Windows Subsystem for Linux](https://learn.microsoft.com/en-us/windows/wsl/) (WSL).

## First-time setup :wrench:

Run `make install-devtools`. This will install...

* [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
* [staticcheck](https://staticcheck.io/)

## Developing :computer:

Running `make precommit` will run all formatters, linters, and the tests. Run this before submitting a PR to ensure the code passes tests and follows our project conventions.

* `make test` will just run the tests
* `make lint` will just run the formatting and linters

## Tests :zap:

Integration tests require an auth token for testing. Set the env var `TEST_AUTH_TOKEN` to
provide it, you can get this from your `~/.momento/credentials` file. The env `TEST_CACHE_NAME` is also required, but for now any string value works.

Then run `make test`.

{{ ossFooter }}
