# Momento client-sdk-go

:warning: Experimental SDK :warning:

Go SDK for Momento is experimental and under active development.
There could be non-backward compatible changes or removal in the future.
Please be aware that you may need to update your source code with the current version of the SDK when its version gets upgraded.

---

<br />

<div align="center">
    <img src="images/gopher.png" alt="Logo" width="200" height="150">
</div>

Go SDK for Momento, a serverless cache that automatically scales without any of the operational overhead required by traditional caching solutions.

<br/>

# Getting Started :running:

## Requirements

- [Go version 0.17.\*](https://go.dev/dl/)
- A Momento Auth Token is required, you can generate one using the [Momento CLI](https://github.com/momentohq/momento-cli)

<br/>

## Installing Momento and Running the Example

Check out our [Go SDK example repo](add_link_here)!

<br />

# Running Tests :zap:

## Requirements

- `TEST_AUTH_TOKEN` - an auth token for testing
- `TEST_CACHE_NAME` - any string value would work

## How to Run Test

`TEST_AUTH_TOKEN=<auth token> TEST_CACHE_NAME=<cache name> go test -v ./momento`
