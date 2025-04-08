package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	proto "github.com/Manan-Rastogi/grpc-sensor-system/proto"
	"github.com/google/uuid"
)

func generateSensorData() *proto.SensorData {
	id := uuid.NewString()
	temp := 50 + rand.Float64()*(100-50)
	timestamp := time.Now().Unix()

	fmt.Println("Generated: ", id)

	return &proto.SensorData{
		Id:          id,
		Temperature: float32(temp),
		Timestamp:   timestamp,
	}
}

func DataProducer(ctx context.Context, SensorDataIn chan<- *proto.SensorData) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer close(SensorDataIn) // clean exit

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Data Generation Ended!!!")
			return
		case <-ticker.C:
			data := generateSensorData()
			SensorDataIn <- data
		}
	}
}

func Processor(inputChan <-chan *proto.SensorData, inputCountChan chan struct{}, outputChan chan<- *proto.ServerResponse, client proto.SensorServiceClient, wg *sync.WaitGroup) {
	for data := range inputChan {
		inputCountChan <- struct{}{} // acquire slot
		wg.Add(1)
		go func(data *proto.SensorData) { // why routine here? --> so defer cancel() is not inside of for loop making the entire loop to finish to cancel the context out..
			defer func() {
				<-inputCountChan
				wg.Done()
			}() // release slot

			

			resp, err := client.SendSensorData(withAuthCtx(), data)
			if err != nil {
				log.Println("Error receiving Data: ", err.Error())
				outputChan <- &proto.ServerResponse{
					Status: "ERROR!",
				}
				return
			}

			outputChan <- resp
		}(data)
	}
}

