package utils

import (
	"context"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// GetRequestID extracts request ID from Lambda context
func GetRequestID(ctx context.Context) string {
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		return lc.AwsRequestID
	}
	return "unknown"
} 