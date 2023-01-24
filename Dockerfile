FROM golang:1.19-alpine
ENV CGO_ENABLED=0
ENV TEST_NAME=default

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /app
WORKDIR /app/incubating
CMD go test -run $TEST_NAME
