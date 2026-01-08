# Running the Example

## Requirements.

- [Go](https://go.dev/dl/)
- You'll need a Momento API key to authenticate with Momento. You can get one from the [Momento Console](https://console.gomomento.com/caches).
- A Momento service endpoint is required. You can find a [list of them here](https://docs.momentohq.com/platform/regions).
- Run `go mod vendor` to install dependencies.

## Running an example.

Each example is a main.go file in its own directory.

To run an example, first provide the necessary environment variables:

```shell
export MOMENTO_API_KEY=<api key>
export MOMENTO_ENDPOINT=<endpoint>
```

Then do `go run <directory>/*.go` for the specific example directory you want to run. 
For example, to run the get/set/delete example...

```
go run scalar-example/*.go
```

## Using SDK in your project

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```bash
go get github.com/momentohq/client-sdk-go
```
