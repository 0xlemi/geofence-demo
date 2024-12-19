package main

import (
	"context"
	"geofence-demo/internal/handler"
	"geofence-demo/internal/metrics"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func main() {
	// Configure AWS with specific region
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		panic(err)
	}

	// Create CloudWatch client
	cwClient := cloudwatch.NewFromConfig(cfg)
	metrics := metrics.New(cwClient)

	h := handler.New(metrics)
	lambda.Start(h.Handle)
}
