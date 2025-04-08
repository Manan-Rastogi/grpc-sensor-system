package main

import (
	"fmt"
	"io"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"google.golang.org/grpc"
)

func (s *SensorServer) LiveSensorChats(stream grpc.BidiStreamingServer[proto.SensorData, proto.ServerResponse]) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("Received from client: %+v\n", data)

		status := "OK"

		if data.GetTemperature() > 80 {
			status = fmt.Sprintf("WARNING: High Temperature for %v", data.GetId())
		}

		err = stream.Send(&proto.ServerResponse{
			Status: status,
		})
		
		if err != nil {
			return err
		}
	}
}
