import {App, Stack} from 'aws-cdk-lib'
import {
    Cluster,
    Compatibility,
    ContainerImage,
    FargateService,
    LogDriver,
    Secret,
    TaskDefinition
} from 'aws-cdk-lib/aws-ecs'
import {resolve} from "path";
import {Vpc} from "aws-cdk-lib/aws-ec2";
import {Platform} from "aws-cdk-lib/aws-ecr-assets";
import {LogGroup} from "aws-cdk-lib/aws-logs";
import * as secretsmanager from "aws-cdk-lib/aws-secretsmanager";

const app = new App();
const stack = new Stack(app, 'pub-sub-demo');
const pubSubSecret = secretsmanager.Secret.fromSecretNameV2(stack, "pub-sub-secret", "pubsub/secret");
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
    environment: {"TEST_TOPIC_NAME": "test-topic", "SERVICE_TO_RUN": "publisher"},
    portMappings: [{containerPort: 3000}],
    logging: LogDriver.awsLogs({streamPrefix: 'publisher', logGroup}),
    secrets: {"TEST_AUTH_TOKEN": Secret.fromSecretsManager(pubSubSecret, "AUTH_TOKEN")}
});
subscriberTaskDefinition.addContainer("subscriber-container", {
    containerName: 'subscriber',
    image: ContainerImage.fromAsset(resolve(__dirname, '..'), {platform: Platform.LINUX_AMD64}),
    environment: {"TEST_TOPIC_NAME": "test-topic", "SERVICE_TO_RUN": "subscriber"},
    logging: LogDriver.awsLogs({streamPrefix: 'subscriber', logGroup}),
    secrets: {"TEST_AUTH_TOKEN": Secret.fromSecretsManager(pubSubSecret, "AUTH_TOKEN")}
});

new FargateService(stack, "publisher-fargate-service", {
   cluster:  publisherCluster,
    taskDefinition: publisherTaskDefinition,
    desiredCount: 0
});
new FargateService(stack, "subscriber-fargate-service", {
    cluster:  subscriberCluster,
    taskDefinition: subscriberTaskDefinition,
    desiredCount: 0
});


app.synth();
