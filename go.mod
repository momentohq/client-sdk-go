module github.com/momentohq/client-sdk-go

go 1.19

require (
	github.com/golang-jwt/jwt/v4 v4.3.0
	github.com/google/uuid v1.6.0
	github.com/onsi/ginkgo/v2 v2.8.1
	github.com/onsi/gomega v1.26.0
	golang.org/x/net v0.21.0
	google.golang.org/grpc v1.61.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace google.golang.org/grpc => github.com/bruuuuuuuce/grpc-go v1.61.0-dev.0.20240306220430-6ad501388d03
