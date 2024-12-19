#!/bin/bash

# AWS Setup Commands Reference
# Order of execution for setting up the Geofence Lambda project

# 1. Set variables
REGION="us-east-2"
FUNCTION_NAME="geofence-demo"
ROLE_NAME="geofence-lambda-role"

# 2. Create IAM Role and attach policies
echo "Creating IAM role..."
aws iam create-role --role-name $ROLE_NAME \
    --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{"Effect": "Allow","Principal": {"Service": "lambda.amazonaws.com"},"Action": "sts:AssumeRole"}]}'

# Wait for role to be created
echo "Waiting for role to propagate..."
sleep 10

# Attach policies
echo "Attaching policies..."
aws iam attach-role-policy --role-name $ROLE_NAME \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

aws iam attach-role-policy --role-name $ROLE_NAME \
    --policy-arn arn:aws:iam::aws:policy/CloudWatchLogsFullAccess

# Add CloudWatch metrics policy
aws iam attach-role-policy --role-name $ROLE_NAME \
    --policy-arn arn:aws:iam::aws:policy/CloudWatchFullAccess

# 3. Get Role ARN
ROLE_ARN=$(aws iam get-role --role-name $ROLE_NAME --query 'Role.Arn' --output text)

# 4. Build the Lambda package
echo "Building Lambda package..."
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go
zip function.zip bootstrap

# 5. Create Lambda function
echo "Creating Lambda function..."
aws lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime provided.al2 \
    --handler bootstrap \
    --role $ROLE_ARN \
    --zip-file fileb://function.zip \
    --memory-size 128 \
    --timeout 30 \
    --environment "Variables={LOG_LEVEL=info}" \
    --region $REGION

# 6. Create CloudWatch Dashboard
echo "Creating CloudWatch dashboard..."
aws cloudwatch put-dashboard --dashboard-name "GeofenceDemo" --dashboard-body '{
    "widgets": [
        {
            "type": "metric",
            "properties": {
                "metrics": [
                    ["GeofenceService", "RequestCount", "Service", "geofence"],
                    ["GeofenceService", "GeofenceHit", "Service", "geofence"],
                    ["GeofenceService", "GeofenceMiss", "Service", "geofence"]
                ],
                "period": 300,
                "stat": "Sum",
                "region": "us-east-2",
                "title": "Geofence Activity"
            }
        }
    ]
}'

# Update commands for future deployments
echo -e "\nFor future updates, use:"
echo "GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go"
echo "zip function.zip bootstrap"
echo "aws lambda update-function-code --function-name $FUNCTION_NAME --zip-file fileb://function.zip" 