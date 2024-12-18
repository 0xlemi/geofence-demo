package handler

import (
	"context"
	"fmt"
	"geofence-demo/internal/geofence"
	"go.uber.org/zap"
)

type Handler struct {
	geoService *geofence.Service
	logger     *zap.Logger
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
	logger, _ := zap.NewProduction()
	return &Handler{
		geoService: geofence.New(),
		logger:     logger,
	}
}

func (h *Handler) Handle(ctx context.Context, req Request) (resp Response, err error) {
	// Defer panic recovery
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.Stack("stack"),
			)
			err = fmt.Errorf("internal error: %v", r)
			resp = Response{
				Message: "Internal server error",
			}
		}
	}()

	// Rest of the handler code...
	return h.handleRequest(ctx, req)
}

// Move existing handler logic to new method
func (h *Handler) handleRequest(ctx context.Context, req Request) (Response, error) {
	h.logger.Info("processing request",
		zap.String("device_id", req.DeviceID),
		zap.Float64("latitude", req.Lat),
		zap.Float64("longitude", req.Lng),
		zap.String("timestamp", req.Timestamp),
	)
	// Validate request
	if req.DeviceID == "" {
		return Response{}, &geofence.ValidationError{
			Field: "device_id",
			Value: req.DeviceID,
			Msg:   "is required",
		}
	}

	// Check coordinates
	if err := h.validateCoordinates(req.Lat, req.Lng); err != nil {
		return Response{}, fmt.Errorf("invalid coordinates: %w", err)
	}

	// Check timestamp
	if req.Timestamp == "" {
		return Response{}, fmt.Errorf("processing failed: %w", 
			&geofence.ValidationError{
				Field: "timestamp",
				Value: req.Timestamp,
				Msg:   "is required",
			})
	}

	// TODO: Add geofence check logic
	return Response{}, nil
}

func (h *Handler) validateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return &geofence.ValidationError{
			Field: "latitude",
			Value: lat,
			Msg:   "must be between -90 and 90",
		}
	}
	if lng < -180 || lng > 180 {
		return &geofence.ValidationError{
			Field: "longitude",
			Value: lng,
			Msg:   "must be between -180 and 180",
		}
	}
	return nil
}
