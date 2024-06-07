# Running the Example

## Requirements.

- [Go](https://go.dev/dl/)
- A Momento Auth Token; you can generate one using the [Momento CLI](https://github.com/momentohq/momento-cli).
- Run `go mod vendor` to install dependencies.

## Running an example.

Each example is a main.go file in its own directory.

To run an example, provide your Momento Auth Token as the MOMENTO_API_KEY environment variable and `go run <directory>/*.go` for the specific exmaple directory. For example, to run the get/set/delete example...

```
MOMENTO_API_KEY=<YOUR_TOKEN> go run scalar-example/*.go
```

## Using SDK in your project

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```bash
go get github.com/momentohq/client-sdk-go
```
