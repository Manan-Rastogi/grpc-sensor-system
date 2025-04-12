package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func BidirectionalSoln(client proto.SensorServiceClient) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "Bearer super-secret-token",
	}))
	// start the stream
	stream, err := client.LiveSensorChats(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// routine to send data
	go func(client proto.SensorServiceClient) {
		for range 10 {
			data := &proto.SensorData{
				Id:          uuid.NewString(),
				Temperature: float32(rand.Intn(100)),
				Timestamp:   time.Now().Unix(),
			}

			err := stream.Send(data)
			if err != nil {
				return
			}

			time.Sleep(1 * time.Second)
		}

		stream.CloseSend()
	}(client)

	// Receiving Data
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf("data: %+v\n", data)
	}
}
