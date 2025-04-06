package main

import (
	"context"
	"log"
	"sync"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
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

	ProducerConsumerSoln(client)
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