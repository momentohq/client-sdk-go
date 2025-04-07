# Compression Middleware Example

This example implements a middleware that uses zstd on a several different read and write requests: Get, Set, GetWithHash, and SetIfAbsentOrHashEqual.

To run the example:

```bash
MOMENTO_API_KEY=<YOUR_TOKEN> go run compression/*.go 
```