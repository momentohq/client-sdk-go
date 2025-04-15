# Compression Middleware Example

This example shows you how to enable gzip compression middleware for a subset of scalar read and write requests: Get, Set, GetWithHash, and SetIfAbsentOrHashEqual.

To run the example:

```bash
MOMENTO_API_KEY=<YOUR_TOKEN> go run compression/main.go 
```