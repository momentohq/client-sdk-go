<div align="center">
    <h1 align="center">Client SDK Go</h1>
    <img src="images/gopher.png" alt="Logo" width="200" height="150">
</div>

# Prerequisites

- Go version 0.17.\*

# Tests

## Requirements

- `TEST_AUTH_TOKEN` - an auth token for testing
- `TEST_CACHE_NAME` - any string value would work

## How to Run Test

`TEST_AUTH_TOKEN=<auth token> TEST_CACHE_NAME=<cache name> go test -v ./momento`
