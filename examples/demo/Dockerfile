FROM golang:1.19-alpine
ENV CGO_ENABLED=0

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./ /app
RUN go build
CMD ./demo
