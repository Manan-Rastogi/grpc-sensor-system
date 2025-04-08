package main

import (
	"context"
	"fmt"
	"io"
	"testing"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

// Mocking grpc.ClientStreamingServer[proto.SensorData, proto.ServerResponse] for our test

type mockClientStream struct {
	recvIndex int
	sentResp  *proto.ServerResponse
	dataBatch []*proto.SensorData
}

func (m *mockClientStream) Recv() (*proto.SensorData, error) {
	if m.recvIndex >= len(m.dataBatch) {
		return nil, io.EOF
	}

	d := m.dataBatch[m.recvIndex]
	m.recvIndex++
	return d, nil
}

func (m *mockClientStream) SendAndClose(res *proto.ServerResponse) error {
	m.sentResp = res
	return nil
}

// Required to satisfy the full gRPC stream interface (even if unused)
func (m *mockClientStream) Context() context.Context     { return context.Background() }
func (m *mockClientStream) SendHeader(metadata.MD) error { return nil }
func (m *mockClientStream) SetHeader(metadata.MD) error  { return nil }
func (m *mockClientStream) SetTrailer(metadata.MD)       {}
func (m *mockClientStream) SendMsg(interface{}) error    { return nil }
func (m *mockClientStream) RecvMsg(interface{}) error    { return nil }

func TestUploadSensorBatch(t *testing.T) {
	mockData := []*proto.SensorData{
		{Id: "1", Temperature: 78},
		{Id: "2", Temperature: 90},
		{Id: "3", Temperature: 55},
	}

	mockStream := &mockClientStream{
		dataBatch: mockData,
	}

	server := &SensorServer{}

	err := server.UploadSensorBatch(mockStream)
	require.NoError(t, err)
	require.NotNil(t, mockStream.sentResp)
	require.Equal(t, fmt.Sprintf("Received %d sensor records.", len(mockData)), mockStream.sentResp.Status)
}
