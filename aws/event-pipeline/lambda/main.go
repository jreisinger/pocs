package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ddbClient *dynamodb.Client
	tableName string
	logger    *slog.Logger
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("load aws config", "err", err)
		panic(err)
	}
	ddbClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("TABLE_NAME")
}

func handler(ctx context.Context, event events.EventBridgeEvent) error {
	start := time.Now()
	log := logger.With(
		"eventId", event.ID,
		"source", event.Source,
		"detailType", event.DetailType,
	)
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		log = log.With("requestId", lc.AwsRequestID)
	}

	log.Info("received event")

	detailJSON, err := json.Marshal(event.Detail)
	if err != nil {
		log.Error("marshal detail", "err", err)
		return fmt.Errorf("marshal detail: %w", err)
	}

	item := map[string]types.AttributeValue{
		"eventId":    &types.AttributeValueMemberS{Value: event.ID},
		"source":     &types.AttributeValueMemberS{Value: event.Source},
		"detailType": &types.AttributeValueMemberS{Value: event.DetailType},
		"detail":     &types.AttributeValueMemberS{Value: string(detailJSON)},
		"receivedAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
	}

	_, err = ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Error("put item", "err", err, "table", tableName)
		return fmt.Errorf("put item: %w", err)
	}

	log.Info("stored event", "durationMs", time.Since(start).Milliseconds())
	return nil
}

func main() {
	lambda.Start(handler)
}
