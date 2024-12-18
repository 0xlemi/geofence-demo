package main

import (
	"geofence-demo/internal/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	h := handler.New()
	lambda.Start(h.Handle)
}
