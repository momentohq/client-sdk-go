# Compression Middleware Example

This example implements a middleware that uses zstd on Get and Set requests. 
Additional methods could be added to the middleware, but we showcase only Get and Set for simplicity.

To run the example:

```bash
MOMENTO_API_KEY=<YOUR_TOKEN> go run compression/*.go 
```