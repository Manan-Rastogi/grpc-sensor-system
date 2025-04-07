package main

import (
	"context"
	"log"
	"math/rand/v2"
	"net"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
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

		time.Sleep(1*time.Second)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen: ", err.Error())
	}

	grpcServer := grpc.NewServer()

	proto.RegisterSensorServiceServer(grpcServer, &SensorServer{})

	log.Println("Server Up and Running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err.Error())
	}
}
