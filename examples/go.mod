module github.com/momentohq/go-example

go 1.19

require (
	// the hrtime and hdrhistogram-go modules are not required to use momento, but
	// they are used in the loadgen example
	github.com/HdrHistogram/hdrhistogram-go v1.1.2
	github.com/google/uuid v1.6.0
	github.com/loov/hrtime v1.0.3
	github.com/momentohq/client-sdk-go v1.20.0

	// logrus is not required to use momento, but it is used in the logging-example
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/aws/aws-cdk-go/awscdk/v2 v2.136.0
	github.com/aws/constructs-go/constructs/v10 v10.3.0
	github.com/aws/jsii-runtime-go v1.97.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/cdklabs/awscdk-asset-awscli-go/awscliv1/v2 v2.2.202 // indirect
	github.com/cdklabs/awscdk-asset-kubectl-go/kubectlv20/v2 v2.1.2 // indirect
	github.com/cdklabs/awscdk-asset-node-proxy-agent-go/nodeproxyagentv6/v2 v2.0.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc v1.63.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/momentohq/client-sdk-go => ../../client-sdk-go
