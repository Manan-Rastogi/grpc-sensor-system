package main

import (
	"context"
	"log"
	"net"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"google.golang.org/grpc"
)

// Implementing SensorServiceServer interface
type SensorServer struct {
	proto.UnimplementedSensorServiceServer
}

func (s *SensorServer) SendSensorData(ctx context.Context, data *proto.SensorData) (*proto.ServerResponse, error) {
	log.Printf("Received Data from Sensors: %+v, %+v, %+v", data.Id, data.Temperature, data.Timestamp)

	// Here we can process data.....
	////////////////////////////////

	time.Sleep(4*time.Second)

	return &proto.ServerResponse{Status: "Received!"}, nil
}

func main(){
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


