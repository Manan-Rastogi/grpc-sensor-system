package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Connect with gRPC servers
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Unable to Connect to server: ", err.Error())
	}

	defer conn.Close()

	client := proto.NewSensorServiceClient(conn)

	// ProducerConsumerSoln(client)

	// ServerStreamingSoln(client)

	ClientStreamigSoln(client)
}

func ClientStreamigSoln(client proto.SensorServiceClient) {
	var batch []*proto.SensorData
	for i := 0; i < 10; i++ {
		batch = append(batch, &proto.SensorData{
			Id:          uuid.New().String(),
			Temperature: rand.Float32()*50 + 20,
			Timestamp:   time.Now().Unix(),
		})
	}

	UploadSensorBatch(client, batch)
}

func UploadSensorBatch(client proto.SensorServiceClient, dataList []*proto.SensorData) {
	stream, err := client.UploadSensorBatch(context.Background())
	if err != nil {
		log.Fatalf("Error starting client stream: %v", err)
	}

	for _, data := range dataList {
		log.Printf("Sending data ID: %s\n", data.Id)
		if err := stream.Send(data); err != nil {
			log.Fatalf("Error sending data: %v", err)
		}
	}

	// Close the stream and receive response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	log.Println("Upload finished with response:", resp.Status)
}

func ServerStreamingSoln(client proto.SensorServiceClient) {
	req := &proto.SensorRequest{
		SensorId: "sensor_112",
	}

	stream, err := client.GetSensorDataStream(context.Background(), req)
	if err != nil {
		log.Fatalf("Error starting stream: %v", err)
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error receiving data: %v", err.Error())
		}

		log.Printf("Streamed Sensor Data: %+v", data)
	}
}

func ProducerConsumerSoln(client proto.SensorServiceClient) {
	// 2. Challange Soln inplemented...
	wg := &sync.WaitGroup{}
	inputCountChan := make(chan struct{}, 5)
	inputChan := make(chan *proto.SensorData, 10)
	outputChan := make(chan *proto.ServerResponse)

	// Produce Data
	ctx_1min, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	go DataProducer(ctx_1min, inputChan) // creating data till certain time

	// Process Data
	go func(inputChan <-chan *proto.SensorData, inputCountChan chan struct{}, outputChan chan<- *proto.ServerResponse, client proto.SensorServiceClient, wg *sync.WaitGroup) {
		Processor(inputChan, inputCountChan, outputChan, client, wg)

		wg.Wait()
		close(outputChan)
	}(inputChan, inputCountChan, outputChan, client, wg)

	// Read Data
	for data := range outputChan {
		log.Println("Received Data: ", data)
	}
	close(inputCountChan)
}
