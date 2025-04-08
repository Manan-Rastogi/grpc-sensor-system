package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/Manan-Rastogi/grpc-sensor-system/server/interceptors"
	"google.golang.org/grpc"
)

// Implementing SensorServiceServer interface
type SensorServer struct {
	proto.UnimplementedSensorServiceServer
}

// Unary: 1 Req-> 1 Res
func (s *SensorServer) SendSensorData(ctx context.Context, data *proto.SensorData) (*proto.ServerResponse, error) {
	log.Printf("Received Data from Sensors: %+v, %+v, %+v", data.Id, data.Temperature, data.Timestamp)

	// Here we can process data.....
	////////////////////////////////

	time.Sleep(4 * time.Second)

	return &proto.ServerResponse{Status: "Received!"}, nil
}

// Server Streaming: 1 Req -> Stream Response
func (s *SensorServer) GetSensorDataStream(req *proto.SensorRequest, stream grpc.ServerStreamingServer[proto.SensorData]) error {
	sensorId := req.GetSensorId()
	log.Println("Streaming data for:", sensorId)

	for range 10 {
		data := &proto.SensorData{
			Id:          sensorId,
			Temperature: float32(50 + rand.Float64()*(100-50)),
			Timestamp:   time.Now().Unix(),
		}

		if err := stream.Send(data); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func (s *SensorServer) UploadSensorBatch(stream grpc.ClientStreamingServer[proto.SensorData, proto.ServerResponse]) error {
	count := 0

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			// All data received
			return stream.SendAndClose(&proto.ServerResponse{
				Status: fmt.Sprintf("Received %d sensor records.", count),
			})
		}
		if err != nil {
			log.Println("Error receiving stream:", err)
			return err
		}

		log.Printf("Received SensorData: ID=%s, Temp=%.2f\n", data.Id, data.Temperature)
		count++
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen: ", err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryAuthInterceptors),
		grpc.StreamInterceptor(interceptors.StreamAuthInterceptor),
	)

	proto.RegisterSensorServiceServer(grpcServer, &SensorServer{})

	log.Println("Server Up and Running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err.Error())
	}
}
