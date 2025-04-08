package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// go test -v -cover ./server -run ^TestIntegration_SendSensorData$
func startTestGRPCServer(t *testing.T) (proto.SensorServiceClient, func()) {
	lis, err := net.Listen("tcp", ":0") // Spin up the server for test

	require.NoError(t, err)

	grpcServer := grpc.NewServer()

	proto.RegisterSensorServiceServer(grpcServer, &SensorServer{})

	// run server in background
	go func() {
		err := grpcServer.Serve(lis)
		require.NoError(t, err)
	}()

	// Dial Server as a REAL Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)

	client := proto.NewSensorServiceClient(conn)

	return client, func() {
		grpcServer.Stop()
		conn.Close()
	}
}

// Unary Test
func TestIntegration_SendSensorData(t *testing.T) {
	client, cleanup := startTestGRPCServer(t)
	defer cleanup()

	req := &proto.SensorData{
		Id:          "int-test-1",
		Temperature: 77,
		Timestamp:   time.Now().Unix(),
	}

	resp, err := client.SendSensorData(context.Background(), req)
	require.NoError(t, err)
	require.Contains(t, resp.GetStatus(), "Received")
}

func TestIntegration_GetSensorDataStream(t *testing.T) {
	client, cleanup := startTestGRPCServer(t)
	defer cleanup()

	stream, err := client.GetSensorDataStream(context.Background(), &proto.SensorRequest{
		SensorId: "test-stream-id",
	})
	require.NoError(t, err)

	count := 0
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Equal(t, "test-stream-id", data.GetId())
		count++
	}

	require.Equal(t, 10, count, "Expected 10 data points to be streamed")
}

func TestIntegration_UploadSensorBatch(t *testing.T) {
	client, cleanup := startTestGRPCServer(t)
	defer cleanup()

	stream, err := client.UploadSensorBatch(context.Background())
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		err := stream.Send(&proto.SensorData{
			Id:          fmt.Sprintf("sensor-%d", i),
			Temperature: 60 + float32(i),
			Timestamp:   time.Now().Unix(),
		})
		require.NoError(t, err)
	}

	resp, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.Contains(t, resp.GetStatus(), "Received 5")
}

func TestIntegration_LiveSensorChats(t *testing.T) {
	client, cleanup := startTestGRPCServer(t)
	defer cleanup()

	stream, err := client.LiveSensorChats(context.Background())
	require.NoError(t, err)

	inputs := []*proto.SensorData{
		{Id: "x1", Temperature: 65, Timestamp: time.Now().Unix()},
		{Id: "x2", Temperature: 82, Timestamp: time.Now().Unix()},
		{Id: "x3", Temperature: 99, Timestamp: time.Now().Unix()},
	}

	for _, data := range inputs {
		err := stream.Send(data)
		require.NoError(t, err)

		resp, err := stream.Recv()
		require.NoError(t, err)

		if data.Temperature > 80 {
			require.Contains(t, resp.GetStatus(), "WARNING")
		} else {
			require.Equal(t, "OK", resp.GetStatus())
		}
	}

	_ = stream.CloseSend()
}
