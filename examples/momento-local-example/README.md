# Momento-local example

This example shows you how to connect to a momento-local server for testing Momento without making calls to the live service.

1. Start a [momento-local Docker container](https://hub.docker.com/r/gomomento/momento-local): 

    ```bash
    docker run -p 8080:8080 gomomento/momento-local
    ```

2. Run the example (no Momento API key necessary):

    ```bash
    go run momento-local-example/main.go
    ```
