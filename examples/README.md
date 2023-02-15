# Running the Example

## Requirements.

- [Go version 1.19](https://go.dev/dl/) or newer.
- A Momento Auth Token; you can generate one using the [Momento CLI](https://github.com/momentohq/momento-cli).
- Run `go mod vendor` to install dependencies.

## Running an example.

Each example is a main.go file in its own directory.

To run an example, provide your Momento Auth Token as the MOMENTO_AUTH_TOKEN environment variable and `go run` the example's main.go. For example, to run the get/set/delete example...

```
MOMENTO_AUTH_TOKEN=<YOUR_TOKEN> go run scalar-example/main.go
```

## Using SDK in your project

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```bash
go get github.com/momentohq/client-sdk-go
```
