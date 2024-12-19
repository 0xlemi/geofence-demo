#!/bin/bash

# Test cases
GUADALAJARA='{"device_id":"test-1","latitude":20.6596,"longitude":-103.3496,"timestamp":"2024-12-19T14:45:00Z"}'
CDMX='{"device_id":"test-2","latitude":19.4326,"longitude":-99.1332,"timestamp":"2024-12-19T14:45:00Z"}'
OCEAN='{"device_id":"test-3","latitude":0.0,"longitude":0.0,"timestamp":"2024-12-19T14:45:00Z"}'

function test_location() {
    local payload=$1
    local name=$2
    echo "Testing $name location..."
    echo "Payload: $payload"
    
    timeout 2s aws lambda invoke \
        --function-name geofence-demo \
        --cli-binary-format raw-in-base64-out \
        --payload "$payload" \
        --region us-east-2 \
        /dev/stdout || true

    echo ""
    echo "----------------------------------------"
}

echo "Running geofence tests..."
echo "----------------------------------------"

test_location "$GUADALAJARA" "Guadalajara (inside fence-1)"
test_location "$CDMX" "CDMX (inside fence-2)"
test_location "$OCEAN" "Ocean (outside all fences)"

echo "All tests completed!" 