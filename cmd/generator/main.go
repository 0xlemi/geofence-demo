package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"geofence-demo/internal/handler"
)

func main() {
	// Generate test points around each fence
	points := []handler.Request{
		// Inside Guadalajara fence
		{
			DeviceID:  "dev1",
			Lat:      20.6597 + (rand.Float64() * 0.01),
			Lng:      -103.3496 + (rand.Float64() * 0.01),
			Timestamp: time.Now().Format(time.RFC3339),
		},
		// Inside CDMX fence
		{
			DeviceID:  "dev2",
			Lat:      19.4326 + (rand.Float64() * 0.01),
			Lng:      -99.1332 + (rand.Float64() * 0.01),
			Timestamp: time.Now().Format(time.RFC3339),
		},
		// Outside any fence
		{
			DeviceID:  "dev3",
			Lat:      25.0000,
			Lng:      -100.0000,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	// Print as JSON
	for _, p := range points {
		json, _ := json.MarshalIndent(p, "", "  ")
		fmt.Println(string(json))
	}
} 