{{ ossHeader }}

# Welcome to client-sdk-go contributing guide :wave:

Thank you for taking your time to contribute to our Go SDK!
This guide will provide you information to start your own development and testing.
Happy coding :dancer:

## Submitting

If you've found a bug, or have a suggestion, please [open an issue in our project](https://github.com/momentohq/client-sdk-go/issues).

If you want to submit a change, please [submit a pull request to our project](https://github.com/momentohq/client-sdk-go/pulls). Use the normal [Github pull request process](https://docs.github.com/en/pull-requests). Please run `make precommit` before submitting your pull request; see below for more information.

## Minimum Go version

Our minimum supported Go version is currently `1.18`. You can download it from [go.dev](https://go.dev/).

## Requirements

To make development easier, we provide a [Makefile](https://golangdocs.com/makefiles-golang) to do common development tasks. If you're on Windows, you can get `make` by installing [Windows Subsystem for Linux](https://learn.microsoft.com/en-us/windows/wsl/) (WSL).

## First-time setup :wrench:

Run `make install-devtools`. This will install...

- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [staticcheck](https://staticcheck.io/)
- [gingko](https://onsi.github.io/ginkgo/) for testing

As a post-setup test, run `ginkgo` from the command line. If the program is not found, try adding go to the path as follows:

`export PATH=$PATH:$(go env GOPATH)/bin`

And retry running `ginkgo`. If that works, then add the above line to your shell profile.

## Developing :computer:

Running `make precommit` will run all formatters, linters, and the tests. Run this before submitting a PR to ensure the code passes tests and follows our project conventions.

- `make test` will just run the tests
- `make lint` will just run the formatting and linters

## Tests :zap:

We use [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/) to write our tests.

Integration tests require an auth token for testing. Set the env var `MOMENTO_API_KEY` to provide it, you can get this from your `~/.momento/credentials` file.

Then run `make test`.

{{ ossFooter }}
