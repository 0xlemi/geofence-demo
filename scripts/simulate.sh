#!/bin/bash

# Exit on any error
set -e

echo "Building simulator..."
mkdir -p bin
go build -o bin/simulator cmd/simulator/main.go

echo "Starting simulation..."
echo "Press Ctrl+C to stop"
echo "----------------------------------------"

./bin/simulator

echo "Simulation stopped" 