import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as secrets from 'aws-cdk-lib/aws-secretsmanager';
import * as go from '@aws-cdk/aws-lambda-go-alpha';

export class MomentoLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const momentoAuthTokenParam = new cdk.CfnParameter(this, 'MomentoApiKey', {
      type: 'String',
      description: 'The Momento API key that will be used to read from the cache.',
      noEcho: true,
    });

    const authTokenSecret = new secrets.Secret(this, 'MomentoSimpleGetApiKey', {
      secretName: 'MomentoSimpleGetApiKey',
      secretStringValue: new cdk.SecretValue(momentoAuthTokenParam.valueAsString),
    });

    const momentoEndpointParam = new cdk.CfnParameter(this, 'MomentoEndpoint', {
      type: 'String',
      description: 'The Momento service endpoint to connect to.',
      noEcho: true,
    });

    const endpointSecret = new secrets.Secret(this, 'MomentoSimpleGetEndpoint', {
      secretName: 'MomentoSimpleGetEndpoint',
      secretStringValue: new cdk.SecretValue(momentoEndpointParam.valueAsString),
    });

    const getLambda = new go.GoFunction(this, 'MomentoLambdaExample', {
      functionName: 'MomentoLambdaExample',
      runtime: lambda.Runtime.PROVIDED_AL2,
      entry: path.join(__dirname, '../../lambda'),
      timeout: cdk.Duration.seconds(30),
      memorySize: 128,
      environment: {
        MOMENTO_API_KEY_SECRET_NAME: authTokenSecret.secretName,
        MOMENTO_ENDPOINT_SECRET_NAME: endpointSecret.secretName,
      },
    });

    authTokenSecret.grantRead(getLambda);
    endpointSecret.grantRead(getLambda);
  }
}
