module github.com/momentohq/go-example

require (
	// the hrtime and hdrhistogram-go modules are not required to use momento, but
	// they are used in the loadgen example
	github.com/HdrHistogram/hdrhistogram-go v1.1.2
	github.com/google/uuid v1.3.0
	github.com/loov/hrtime v1.0.3
	github.com/momentohq/client-sdk-go v0.16.0

	// logrus is not required to use momento, but it is used in the logging-example
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/grpc v1.52.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
