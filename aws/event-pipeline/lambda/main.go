package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
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
	ddbClient   *dynamodb.Client
	tableName   string
	totalsTable string
	logger      *slog.Logger
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
	totalsTable = os.Getenv("TOTALS_TABLE_NAME")
}

type order struct {
	CustomerID string  `json:"customerId"`
	OrderID    string  `json:"orderId"`
	Amount     float64 `json:"amount"`
}

func parseOrder(detail json.RawMessage) (order, error) {
	var o order
	if err := json.Unmarshal(detail, &o); err != nil {
		return o, fmt.Errorf("unmarshal detail: %w", err)
	}
	if o.CustomerID == "" || o.OrderID == "" {
		return o, errors.New("customerId and orderId are required")
	}
	if o.Amount <= 0 {
		return o, errors.New("amount must be positive")
	}
	return o, nil
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

	// Invalid orders fail the invocation: Lambda retries, then sends to the DLQ.
	o, err := parseOrder(event.Detail)
	if err != nil {
		log.Error("invalid order", "err", err)
		return err
	}
	log = log.With("customerId", o.CustomerID, "orderId", o.OrderID)

	// Store the event and bump the customer's total atomically. The condition
	// makes redelivered events (EventBridge is at-least-once) a no-op.
	_, err = ddbClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Put: &types.Put{
				TableName: aws.String(tableName),
				Item: map[string]types.AttributeValue{
					"eventId":    &types.AttributeValueMemberS{Value: event.ID},
					"source":     &types.AttributeValueMemberS{Value: event.Source},
					"detailType": &types.AttributeValueMemberS{Value: event.DetailType},
					"detail":     &types.AttributeValueMemberS{Value: string(event.Detail)},
					"receivedAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
				},
				ConditionExpression: aws.String("attribute_not_exists(eventId)"),
			}},
			{Update: &types.Update{
				TableName:        aws.String(totalsTable),
				Key:              map[string]types.AttributeValue{"customerId": &types.AttributeValueMemberS{Value: o.CustomerID}},
				UpdateExpression: aws.String("ADD #total :amount, orderCount :one"),
				ExpressionAttributeNames: map[string]string{
					"#total": "total", // reserved word
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":amount": &types.AttributeValueMemberN{Value: strconv.FormatFloat(o.Amount, 'f', -1, 64)},
					":one":    &types.AttributeValueMemberN{Value: "1"},
				},
			}},
		},
	})
	var tce *types.TransactionCanceledException
	if errors.As(err, &tce) && len(tce.CancellationReasons) > 0 &&
		aws.ToString(tce.CancellationReasons[0].Code) == "ConditionalCheckFailed" {
		log.Info("duplicate event, skipped")
		return nil
	}
	if err != nil {
		log.Error("transact write", "err", err)
		return fmt.Errorf("transact write: %w", err)
	}

	log.Info("processed order", "amount", o.Amount, "durationMs", time.Since(start).Milliseconds())
	return nil
}

func main() {
	lambda.Start(handler)
}
