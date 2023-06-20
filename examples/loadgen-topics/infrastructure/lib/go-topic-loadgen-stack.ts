import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import {aws_s3} from 'aws-cdk-lib';
import * as iam from "aws-cdk-lib/aws-iam";
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import {Construct} from 'constructs';
import {Platform} from 'aws-cdk-lib/aws-ecr-assets';
import {BucketAccessControl} from 'aws-cdk-lib/aws-s3';

export class GoTopicLoadgenStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const vpc = new ec2.Vpc(this, 'GoTopicLoadgenVpc', {
      maxAzs: 1, // Default is all AZs in region
    });

    const cluster = new ecs.Cluster(this, 'GoTopicLoadgenCluster', {
      vpc: vpc,
    });

    const goTopicLoadgenImage = ecs.ContainerImage.fromAsset(
      path.join(__dirname, "../../docker"), {platform: Platform.LINUX_AMD64}
    );

    const goTopicLoadgenTaskRole = new iam.Role(this, "go-topic-loadgen-task-role", {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
    });

    const goTopicLoadgenS3Bucket = new aws_s3.Bucket(this, 'go-topic-loadgen-bucket', {
      bucketName: 'go-topic-loadgen-bucket',
      blockPublicAccess: aws_s3.BlockPublicAccess.BLOCK_ALL
    });
    goTopicLoadgenS3Bucket.grantWrite(goTopicLoadgenTaskRole);

    const goTopicLoadgenDefinition = new ecs.FargateTaskDefinition(
      this,
      "go-topic-loadgen-task-def",
      {
        memoryLimitMiB: 4096,
        taskRole: goTopicLoadgenTaskRole,
        cpu: 2048,
        runtimePlatform: {
          operatingSystemFamily: ecs.OperatingSystemFamily.LINUX,
          cpuArchitecture: ecs.CpuArchitecture.X86_64,
        },
        family: "go-topic-loadgen-testing",
      },
    );

    const goTopicLoadgenContainerDefinition = new ecs.ContainerDefinition(
      this,
      "go-topic-loadgen-container-definition",
      {
        taskDefinition: goTopicLoadgenDefinition,
        image: goTopicLoadgenImage,
        logging: ecs.LogDriver.awsLogs({
          streamPrefix: "go-topic-loadgen-log-driver",
        }),
      });
  }
}
