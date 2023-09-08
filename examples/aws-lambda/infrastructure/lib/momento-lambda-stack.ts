import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as secrets from 'aws-cdk-lib/aws-secretsmanager';
import * as go from '@aws-cdk/aws-lambda-go-alpha';

export class MomentoLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const momentoAuthTokenParam = new cdk.CfnParameter(this, 'MomentoAuthToken', {
      type: 'String',
      description: 'The Momento Auth Token that will be used to read from the cache.',
      noEcho: true,
    });

    const authTokenSecret = new secrets.Secret(this, 'MomentoSimpleGetAuthToken', {
      secretName: 'MomentoSimpleGetAuthToken',
      secretStringValue: new cdk.SecretValue(momentoAuthTokenParam.valueAsString),
    });

    const getLambda = new go.GoFunction(this, 'MomentoLambdaExample', {
      functionName: 'MomentoLambdaExample',
      runtime: lambda.Runtime.GO_1_X,
      entry: path.join(__dirname, '../../lambda'),
      timeout: cdk.Duration.seconds(30),
      memorySize: 128,
      environment: {
        MOMENTO_API_KEY_SECRET_NAME: authTokenSecret.secretName,
      },
    });

    authTokenSecret.grantRead(getLambda);
  }
}
