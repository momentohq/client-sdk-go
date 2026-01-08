# Compression Middleware Example

## Requirements.

- [Go](https://go.dev/dl/)
- A Momento API key is required, you can generate one using the [Momento Console](https://console.gomomento.com/api-keys)
- A Momento service endpoint is required. You can find a [list of them here](https://docs.momentohq.com/platform/regions)

## Usage

This example shows you how to enable gzip compression middleware on scalar read and write requests.

To run the example:

```bash
MOMENTO_API_KEY=<YOUR_TOKEN> MOMENTO_ENDPOINT=<endpoint> go run compression/main.go 
```