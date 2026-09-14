# Event-Driven Processing System

A simple event-driven system built with AWS CDK. Events are published to an
Amazon EventBridge bus, which triggers a Lambda function that stores each
event in a DynamoDB table.

```mermaid
flowchart LR
    P[Event Producer] -->|PutEvents| EB[(Amazon EventBridge<br/>Event Bus)]
    EB -->|rule match<br/>invokes| L[Lambda: Event Processor]
    L -->|PutItem| DB[(DynamoDB Table)]
    EB -.->|failed invocations| DLQ[(Dead-Letter Queue)]
```

## Components

- **EventBridge Bus** — receives events from producers.
- **EventBridge Rule** — matches events and invokes the Lambda.
- **Lambda: Event Processor** — writes each incoming event to DynamoDB.
- **DynamoDB Table** — stores processed events.
- **Dead-Letter Queue** — captures events the Lambda fails to process.

## Ordering

EventBridge does not guarantee event ordering — events can be delivered to
the Lambda out of order, even from the same producer. If ordering matters,
include a sequence number or timestamp in the event `Detail` and handle
reordering in the consumer.

## EventBridge vs SNS vs SQS

| | EventBridge | SNS | SQS |
|---|---|---|---|
| Model | Event bus (content-based routing) | Pub/sub (topic-based) | Message queue |
| Routing | Rules match on event content | Fixed topic, optional filter policy | N/A — single queue, polled by consumers |
| Fan-out | Yes, to multiple targets per rule | Yes, to multiple subscribers per topic | No — one message consumed once |
| Ordering | Not guaranteed | Not guaranteed (FIFO topics excepted) | Not guaranteed (FIFO queues excepted) |
| Retention/replay | No (unless Archive/Replay enabled) | No | Yes, until consumed (up to 14 days) |
| Typical use case | Routing events between decoupled services/apps based on content | Fanning out a notification to many subscribers (email, SMS, Lambda, SQS) | Buffering/decoupling work between a producer and a single consumer |

## Project Structure

```
infra/    CDK app (TypeScript) — defines the EventBridge bus, rule,
          DynamoDB table, DLQ, and Lambda function
lambda/   Event processor (Go) — writes each EventBridge event to DynamoDB
```

## Deploy

```sh
cd infra
npm install
npx cdk deploy
```

## Test

Send a test event to the bus (replace the bus/table names with the deployed
stack outputs):

```sh
aws events put-events --entries '[{
  "Source": "test",
  "DetailType": "TestEvent",
  "Detail": "{\"foo\":\"bar\"}",
  "EventBusName": "event-pipeline-bus"
}]'

aws dynamodb scan --table-name <table-name>
```

## Destroy

```sh
cd infra
npx cdk destroy
```

This removes every resource the stack created: the EventBridge bus/rule, the
Lambda function and its IAM role, the DynamoDB table, the dead-letter queue,
and the Lambda's CloudWatch log group.

Not removed (shared across all CDK apps in this AWS account/region, managed
outside this stack): the CDK bootstrap resources from `cdk bootstrap` (S3
asset bucket, ECR repo, IAM roles) and the asset objects already uploaded to
that bucket by past deploys.
