FROM golang:1.20

COPY docker_src .
WORKDIR app
RUN go mod tidy

CMD go run ./loadgen-topics
