<head>
  <meta name="Momento Node.js Client Library Documentation" content="Node.js client software development kit for Momento Cache">
</head>
<img src="https://docs.momentohq.com/img/logo.svg" alt="logo" width="400"/>

[![project status](https://momentohq.github.io/standards-and-practices/badges/project-status-official.svg)](https://github.com/momentohq/standards-and-practices/blob/main/docs/momento-on-github.md)
[![project stability](https://momentohq.github.io/standards-and-practices/badges/project-stability-stable.svg)](https://github.com/momentohq/standards-and-practices/blob/main/docs/momento-on-github.md)

<br>

## Simple Get Lambda

This repo contains an example AWS Lambda Function, built using AWS CDK, that sets and gets an item in a Momento cache.

## Prerequisites

- Node version 14 or higher is required
- To get started with Momento you will need a Momento Auth Token. You can get one from the [Momento Console](https://console.gomomento.com). Check out the [getting started](https://docs.momentohq.com/getting-started) guide for more information on obtaining an auth token.

## Deploying the Simple Get Lambda

First make sure to start Docker and install the dependencies in the `lambda` directory, which is where the AWS Lambda Go handler function lives. 
(TODO: correct make instruction)

```bash
cd lambda
go install
```

The source code for the CDK application lives in the `infrastructure` directory and is defined using TypeScript. To build and deploy it you will first need to install the dependencies:

```bash
cd infrastructure
npm install
```

To deploy the CDK app you will need to have [configured your AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-authentication.html#cli-chap-authentication-precedence).

You will also need a cache and superuser token generated from the [Momento Console](https://console.gomomento.com).

Then run:

```
npm run cdk -- deploy --parameters MomentoAuthToken=<YOUR_MOMENTO_AUTH_TOKEN>
```

The lambda does not set up a way to access itself externally, so to run it, you will have to go to MomentoSimpleGet in AWS Lambda and run a test.

The lambda is set up to make get calls for the key 'key' in the cache 'cache' by default. It does not create a cache or write anything to that key. While it still may give useful latency information if it can't find a cache or key, creating them will let you test in a more realistic way.

If you have the [Momento CLI](https://github.com/momentohq/momento-cli) installed, you can create a cache like this:

```commandline
momento cache create cache
```

You can then set a value for the key:

```commandline
momento cache set key value
```

You can also create a cache and key using the [Momento Console](https://console.gomomento.com).

Finally, you can edit [handler.ts](lambda/simple-get/handler.ts) to change the cache and key the lambda looks for.