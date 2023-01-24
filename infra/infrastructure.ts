import {App, Duration, Stack} from 'aws-cdk-lib'
import {Cluster, Compatibility, ContainerImage, FargateService, LogDriver, TaskDefinition} from 'aws-cdk-lib/aws-ecs'
import {resolve} from "path";
import {Vpc} from "aws-cdk-lib/aws-ec2";
import {Platform} from "aws-cdk-lib/aws-ecr-assets";
import {LogGroup} from "aws-cdk-lib/aws-logs";
import {DnsRecordType, PrivateDnsNamespace, RoutingPolicy, Service} from "aws-cdk-lib/aws-servicediscovery";

const app = new App();
const stack = new Stack(app, 'pub-sub-demo');

const vpc = new Vpc(stack, 'pub-sub-vpc', { maxAzs: 2 });
const publisherCluster = new Cluster(stack, 'publisher-cluster', {vpc});
const subscriberCluster = new Cluster(stack, 'subscriber-cluster', {vpc});

const publisherTaskDefinition = new TaskDefinition(stack, "publisher-task", {
    compatibility: Compatibility.FARGATE,
    cpu: "256",
    memoryMiB: "512"
});
const subscriberTaskDefinition = new TaskDefinition(stack, "subscriber-task", {
    compatibility: Compatibility.FARGATE,
    cpu: "256",
    memoryMiB: "512"
});
const logGroup = new LogGroup(stack, 'pub-sub-log-group', {logGroupName: "pubsub"});
publisherTaskDefinition.addContainer("publisher-container", {
    containerName: 'publisher',
    image: ContainerImage.fromAsset(resolve(__dirname, '..'), {platform: Platform.LINUX_AMD64}),
    environment: {"TEST_TOPIC_NAME": "test-topic", "TEST_NAME": "TestBasicHappyPathPublisher"},
    portMappings: [{containerPort: 3000}],
    logging: LogDriver.awsLogs({streamPrefix: 'publisher', logGroup})
});
subscriberTaskDefinition.addContainer("subscriber-container", {
    containerName: 'subscriber',
    image: ContainerImage.fromAsset(resolve(__dirname, '..'), {platform: Platform.LINUX_AMD64}),
    environment: {"TEST_TOPIC_NAME": "test-topic", "TEST_NAME": "TestBasicHappyPathSubscriber"},
    logging: LogDriver.awsLogs({streamPrefix: 'subscriber', logGroup})
});

const cloudMapNamespace = new PrivateDnsNamespace(stack, 'pub-sub-service-discovery-namespace', {
    name: 'pubsub.com',
    vpc: vpc
});
const pubsubMapService = new Service(stack, 'pub-sub-service-discovery', {
    namespace: cloudMapNamespace,
    dnsRecordType: DnsRecordType.A,
    dnsTtl: Duration.seconds(300),
    name: 'pub-sub',
    routingPolicy: RoutingPolicy.WEIGHTED,
    loadBalancer: true,
});
const publisherService = new FargateService(stack, "publisher-fargate-service", {
   cluster:  publisherCluster,
    taskDefinition: publisherTaskDefinition,
    desiredCount: 1
});
const subscriberService = new FargateService(stack, "subscriber-fargate-service", {
    cluster:  subscriberCluster,
    taskDefinition: subscriberTaskDefinition,
    desiredCount: 1
});

publisherService.associateCloudMapService({service: pubsubMapService})
subscriberService.associateCloudMapService({service: pubsubMapService})

app.synth();
