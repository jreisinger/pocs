import * as path from 'path';
import * as cdk from 'aws-cdk-lib/core';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import { GoFunction } from '@aws-cdk/aws-lambda-go-alpha';
import { Construct } from 'constructs';

export class EventPipelineStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const table = new dynamodb.Table(this, 'EventsTable', {
      partitionKey: { name: 'eventId', type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const bus = new events.EventBus(this, 'EventBus', {
      eventBusName: 'event-pipeline-bus',
    });

    const deadLetterQueue = new sqs.Queue(this, 'ProcessorDLQ');

    const logGroup = new logs.LogGroup(this, 'EventProcessorLogGroup', {
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const processorFn = new GoFunction(this, 'EventProcessor', {
      entry: path.join(__dirname, '..', '..', 'lambda'),
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        TABLE_NAME: table.tableName,
      },
      deadLetterQueue,
      logGroup,
    });

    table.grantWriteData(processorFn);

    new events.Rule(this, 'ProcessAllEventsRule', {
      eventBus: bus,
      eventPattern: {
        detailType: events.Match.exists(),
      },
      targets: [
        new targets.LambdaFunction(processorFn, {
          deadLetterQueue,
          retryAttempts: 2,
        }),
      ],
    });
  }
}
