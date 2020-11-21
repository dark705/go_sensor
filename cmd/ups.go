package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dark705/go_sensor/internal/config"
	"github.com/dark705/go_sensor/internal/helper"
	"github.com/dark705/go_sensor/internal/message"
	"github.com/dark705/go_sensor/internal/rmq"
	nut "github.com/robbiet480/go.nut"
)

func main() {
	//config
	config := сonfig.ReadFromFile("../config/config.yaml")

	//create RMQ connect
	rabbitMQ, err := rmq.NewRMQ(config.ServiceSensor.RMQ)
	helper.FailOnError(err, "Fail on start RMQ connect")
	defer func() {
		log.Println("Close RMQ connect")
		err = rabbitMQ.CloseConnect()
		if err != nil {
			log.Println("Error on close RMQ connect")
		}
	}()

	//sensor
	nutClient, err := nut.Connect(config.ServiceSensor.Device.Nut.IP)
	helper.FailOnError(err, "Failed on connect to nut")
	defer func() {
		log.Println("Close NUT connect")
		_, err = nutClient.Disconnect()
		if err != nil {
			log.Println("Error on close NUT connect")
		}
	}()

	ticker := time.NewTicker(time.Second * time.Duration(config.ServiceSensor.Device.Nut.Ticker))

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//send message to open RMQ connect by ticker
	go func() {
		ReadUPSAndSendRmq(nutClient, rabbitMQ)
		for {
			select {
			case <-ticker.C:
				ReadUPSAndSendRmq(nutClient, rabbitMQ)
			}
		}
	}()

	log.Printf("Got OS signal: %v. Exit...\n", <-osSignals)
	ticker.Stop()
}

func ReadUPSAndSendRmq(nutClient nut.Client, rabbitMQ *rmq.RMQ) {
	jsonMessage, err := message.GetMessageUPSAsJSON(nutClient)
	helper.FailOnError(err, "Fail get message ups")
	err = rabbitMQ.Send(jsonMessage)
	helper.FailOnError(err, "Fail on send RMQ message")
	log.Printf("Succes sent RMQ message:%s", jsonMessage)
}
