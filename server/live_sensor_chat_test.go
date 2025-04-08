package main

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type mockBiDirStream struct {
	dataReceiver []*proto.SensorData
	dataSender   []*proto.ServerResponse
	currIndex    int
}

func (m *mockBiDirStream) Recv() (*proto.SensorData, error) {
	if m.currIndex >= len(m.dataReceiver) {
		return nil, io.EOF
	}

	data := m.dataReceiver[m.currIndex]
	m.currIndex++

	return data, nil
}

func (m *mockBiDirStream) Send(res *proto.ServerResponse) error {
	m.dataSender = append(m.dataSender, res)
	return nil
}

func (m *mockBiDirStream) Context() context.Context     { return context.Background() }
func (m *mockBiDirStream) RecvMsg(mm any) error         { return nil }
func (m *mockBiDirStream) SendHeader(metadata.MD) error { return nil }
func (m *mockBiDirStream) SendMsg(mm any) error         { return nil }
func (m *mockBiDirStream) SetHeader(metadata.MD) error  { return nil }
func (m *mockBiDirStream) SetTrailer(metadata.MD)       {}

func TestLiveSensorChats(t *testing.T) {
	streamInputs := []*proto.SensorData{
		{
			Id:          "id1",
			Temperature: 55,
			Timestamp:   time.Now().Unix(),
		}, {
			Id:          "id2",
			Temperature: 70,
			Timestamp:   time.Now().Unix(),
		}, {
			Id:          "id3",
			Temperature: 85,
			Timestamp:   time.Now().Unix(),
		}, {
			Id:          "id4",
			Temperature: 85,
			Timestamp:   time.Now().Unix(),
		}, {
			Id:          "id5",
			Temperature: 75,
			Timestamp:   time.Now().Unix(),
		}, {
			Id:          "id6",
			Temperature: 99,
			Timestamp:   time.Now().Unix(),
		},
	}

	server := SensorServer{}
	stream := mockBiDirStream{
		dataReceiver: streamInputs,
	}

	err := server.LiveSensorChats(&stream)

	require.NoError(t, err)
	require.Equal(t, len(stream.dataReceiver), len(stream.dataSender))

	warnCount := 0
	for _, data := range stream.dataSender {
		if strings.Contains(data.GetStatus(), "WARNING") {
			warnCount++
		}
	}

	require.Equal(t, 3, warnCount, "Expected 3 warnings for temperatures > 80")
}
