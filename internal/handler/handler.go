package handler

import (
	"context"
	"fmt"
	"geofence-demo/internal/geofence"
)

type Handler struct {
	geoService *geofence.Service
}

type Request struct {
	DeviceID  string  `json:"device_id"`
	Lat       float64 `json:"latitude"`
	Lng       float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type Response struct {
	InGeofence bool   `json:"in_geofence"`
	FenceID    string `json:"fence_id"`
	Message    string `json:"message"`
}

func New() *Handler {
	return &Handler{
		geoService: geofence.New(),
	}
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	// Validate request
	if req.DeviceID == "" {
		return Response{}, fmt.Errorf("device_id is required")
	}
	if req.Lat < -90 || req.Lat > 90 {
		return Response{}, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}
	if req.Lng < -180 || req.Lng > 180 {
		return Response{}, fmt.Errorf("invalid longitude: must be between -180 and 180")
	}
	if req.Timestamp == "" {
		return Response{}, fmt.Errorf("timestamp is required")
	}

	// TODO: Add geofence check logic
	return Response{}, nil
}
