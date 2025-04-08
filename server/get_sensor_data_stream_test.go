package main

import (
	"testing"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Test for server stream is not straight forward.. we need to mock the stream interface

// implementing grpc.ServerStreamingServer
type mockSensorDataStream struct {
	grpc.ServerStream                     // interface embedding...
	received          []*proto.SensorData // since we are not sending data to client.. we will check the responses of server stream in this
}

// implementing send(). this will append the data server is sending to received slice
func (m *mockSensorDataStream) Send(data *proto.SensorData) error {
	m.received = append(m.received, data)
	return nil
}

func TestGetSensorDataStream(t *testing.T) {
	sensorID := "test_12121"

	server := &SensorServer{}
	mockStream := &mockSensorDataStream{}

	req := &proto.SensorRequest{
		SensorId: sensorID,
	}

	err := server.GetSensorDataStream(req, mockStream)

	require.NoError(t, err)
	require.Len(t, mockStream.received, 10, "Expected 10 streamed responses")

	for _, data := range mockStream.received {
		require.Equal(t, "test_12121", data.GetId())
		require.GreaterOrEqual(t, data.GetTemperature(), float32(50))
		require.LessOrEqual(t, data.GetTemperature(), float32(100))
		require.Greater(t, data.GetTimestamp(), int64(0))
	}
}
