package geofence

import "fmt"

type ValidationError struct {
    Field string
    Value interface{}
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s=%v %s", e.Field, e.Value, e.Msg)
}

type GeofenceError struct {
    FenceID string
    Msg     string
}

func (e *GeofenceError) Error() string {
    return fmt.Sprintf("geofence error: %s %s", e.FenceID, e.Msg)
} 