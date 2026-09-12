#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib/core';
import { EventPipelineStack } from '../lib/event-pipeline-stack';

const app = new cdk.App();
new EventPipelineStack(app, 'EventPipelineStack', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});
