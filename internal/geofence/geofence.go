package geofence

import "math"

const earthRadiusMeters = 6371000 // Earth's radius in meters

type Geofence struct {
	ID           string  `json:"id"`
	CenterLat    float64 `json:"center_lat"`
	CenterLng    float64 `json:"center_lng"`
	RadiusMeters float64 `json:"radius_meters"`
}

type Service struct {
	fences []Geofence
}

func New() *Service {
	return &Service{
		fences: []Geofence{
			{
				ID:           "fence-1",
				CenterLat:    20.6597, // Guadalajara Centro
				CenterLng:    -103.3496,
				RadiusMeters: 5000,
			},
			{
				ID:           "fence-2",
				CenterLat:    19.4326, // CDMX Zócalo
				CenterLng:    -99.1332,
				RadiusMeters: 3000,
			},
			{
				ID:           "fence-3",
				CenterLat:    22.8905, // Los Cabos Marina
				CenterLng:    -109.9167,
				RadiusMeters: 2000,
			},
		},
	}
}

// calculateDistance returns distance in meters between two points using haversine formula
func (s *Service) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	lat1Rad := toRadians(lat1)
	lng1Rad := toRadians(lng1)
	lat2Rad := toRadians(lat2)
	lng2Rad := toRadians(lng2)

	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadiusMeters * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (s *Service) IsPointInFence(lat, lng float64) (bool, string) {
	for _, fence := range s.fences {
		dist := s.calculateDistance(lat, lng, fence.CenterLat, fence.CenterLng)
		if dist <= fence.RadiusMeters {
			return true, fence.ID
		}
	}
	return false, ""
}
