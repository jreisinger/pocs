import * as cdk from 'aws-cdk-lib/core';
import { Template } from 'aws-cdk-lib/assertions';
import { EventPipelineStack } from '../lib/event-pipeline-stack';

test('creates event bus, DynamoDB table, and Lambda processor', () => {
  const app = new cdk.App();
  const stack = new EventPipelineStack(app, 'MyTestStack');
  const template = Template.fromStack(stack);

  template.hasResourceProperties('AWS::Events::EventBus', {
    Name: 'event-pipeline-bus',
  });

  template.hasResourceProperties('AWS::DynamoDB::Table', {
    KeySchema: [{ AttributeName: 'eventId', KeyType: 'HASH' }],
  });

  template.resourceCountIs('AWS::Lambda::Function', 1);
  template.resourceCountIs('AWS::SQS::Queue', 1);
});
