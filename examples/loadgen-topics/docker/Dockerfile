FROM golang:1.20

COPY docker_src .
WORKDIR app
RUN go mod tidy -go=1.16 && go mod tidy -go=1.17

CMD go run ./loadgen-topics
