package metrics

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/aws"
	"time"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"fmt"
)

type Metrics struct {
	client *cloudwatch.Client
	namespace string
}

func New(client *cloudwatch.Client) *Metrics {
	return &Metrics{
		client: client,
		namespace: "GeofenceService",
	}
}

func (m *Metrics) IncrementRequests(ctx context.Context) error {
	_, err := m.client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(m.namespace),
		MetricData: []types.MetricDatum{
			{
				MetricName: aws.String("RequestCount"),
				Value:     aws.Float64(1.0),
				Unit:      types.StandardUnitCount,
				Timestamp: aws.Time(time.Now()),
				Dimensions: []types.Dimension{
					{
						Name:  aws.String("Service"),
						Value: aws.String("geofence"),
					},
				},
			},
		},
	})
	return err
}

func (m *Metrics) TrackGeofenceHit(ctx context.Context, fenceID string, isInside bool) error {
	metricName := "GeofenceHit"
	if !isInside {
		metricName = "GeofenceMiss"
	}
	
	input := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(m.namespace),
		MetricData: []types.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Value:     aws.Float64(1.0),
				Unit:      types.StandardUnitCount,
				Timestamp: aws.Time(time.Now()),
				Dimensions: []types.Dimension{
					{
						Name:  aws.String("Service"),
						Value: aws.String("geofence"),
					},
				},
			},
		},
	}
	
	_, err := m.client.PutMetricData(ctx, input)
	if err != nil {
		fmt.Printf("Error publishing metric %s: %v\n", metricName, err)
	}
	return err
} 