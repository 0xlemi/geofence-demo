package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type Location struct {
	DeviceID  string  `json:"device_id"`
	Lat       float64 `json:"latitude"`
	Lng       float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

// Fence boundaries
var fences = []struct {
	name      string
	centerLat float64
	centerLng float64
	radius    float64
}{
	{"Guadalajara", 20.6596, -103.3496, 0.1},
	{"CDMX", 19.4326, -99.1332, 0.1},
}

var (
	totalRequests uint64
	deviceCount   = 10
	interval      = 500 * time.Millisecond
)

func main() {
	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting simulation with %d devices sending every %v\n", deviceCount, interval)

	// Start stats printer
	go printStats()

	// Start devices
	for i := 0; i < deviceCount; i++ {
		deviceID := fmt.Sprintf("sim-device-%d", i)
		go simulateDevice(deviceID)
	}

	// Wait for interrupt
	<-sigChan
	fmt.Printf("\nSimulation stopped. Total requests: %d\n", atomic.LoadUint64(&totalRequests))
}

func simulateDevice(deviceID string) {
	// Initial position near a random fence
	fence := fences[rand.Intn(len(fences))]

	for {
		// Randomly decide to be inside or outside fence
		if rand.Float64() < 0.3 { // 30% chance to be outside
			// Way outside the fence
			fence = fences[rand.Intn(len(fences))]
			fence.centerLat += 0.2 // 0.2 degrees is well outside our 0.1 radius
			fence.centerLng += 0.2
		} else {
			// Inside or near the fence
			fence = fences[rand.Intn(len(fences))]
			fence.centerLat += (rand.Float64() - 0.5) * 0.1
			fence.centerLng += (rand.Float64() - 0.5) * 0.1
		}

		loc := Location{
			DeviceID:  deviceID,
			Lat:       fence.centerLat,
			Lng:       fence.centerLng,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		sendLocation(loc)
		atomic.AddUint64(&totalRequests, 1)
		time.Sleep(interval)
	}
}

func sendLocation(loc Location) {
	data, _ := json.Marshal(loc)
	cmd := fmt.Sprintf("aws lambda invoke --function-name geofence-demo --payload '%s' --cli-binary-format raw-in-base64-out --region us-east-2 /dev/null", string(data))
	
	if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
		fmt.Printf("Error sending location for %s: %v\n", loc.DeviceID, err)
	}
}

func printStats() {
	ticker := time.NewTicker(1 * time.Second)
	var lastCount uint64

	for range ticker.C {
		currentCount := atomic.LoadUint64(&totalRequests)
		rps := currentCount - lastCount
		fmt.Printf("\rRequests: %d (%.2f/sec)", currentCount, float64(rps))
		lastCount = currentCount
	}
}
