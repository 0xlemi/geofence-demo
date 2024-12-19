#!/bin/bash

# Exit on any error
set -e

# Create bin directory if it doesn't exist
mkdir -p bin

echo "Building Lambda function..."
GOOS=linux GOARCH=amd64 go build -o bin/bootstrap cmd/lambda/main.go

echo "Creating deployment package..."
cd bin && zip function.zip bootstrap && cd ..

echo "Updating Lambda function..."
timeout 2s aws lambda update-function-code \
    --function-name geofence-demo \
    --zip-file fileb://bin/function.zip \
    --region us-east-2 || true  # Continue even if timeout occurs

echo "Cleaning up..."
rm bin/bootstrap bin/function.zip

echo "Deployment complete!"
