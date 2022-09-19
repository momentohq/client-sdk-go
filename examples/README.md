## Running the Example

- [Go version 1.18.\*](https://go.dev/dl/) is required
- A Momento Auth Token is required, you can generate one using the [Momento CLI](https://github.com/momentohq/momento-cli)

```bash
go mod vendor
MOMENTO_AUTH_TOKEN=<YOUR_TOKEN> go run main.go
```

Code example can be found [here](main.go)!

<br />

## Using SDK in your project

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```bash
go get github.com/momentohq/client-sdk-go
```
