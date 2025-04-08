package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/stretchr/testify/require"
)

func TestSendSensorData(t *testing.T) {
	server := &SensorServer{}

	req := &proto.SensorData{
		Id:          "test123",
		Temperature: 85,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()

	resp, err := server.SendSensorData(ctx, req)

	require.NoError(t, err, fmt.Sprintf("No Error Expected but got %v", err))

	// Assert that actualString contains expectedSubstring. If not, fail the test and show the optional error message.
	require.Contains(t, resp.GetStatus(), "Received!", fmt.Sprintf("Data was not received-> %+v", resp))
	
}
