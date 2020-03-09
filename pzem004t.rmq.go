package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/be-ys/pzem-004t-v3/pzem"
	"github.com/dark705/go_sensor/internal/config"
	"github.com/dark705/go_sensor/internal/helper"
	"github.com/dark705/go_sensor/internal/message"
	"github.com/dark705/go_sensor/internal/rmq"
)

func main() {
	//config
	config := сonfig.ReadFromFile("config/config.yaml")

	//create RMQ connect
	r, err := rmq.NewRMQ(config.ServiceSensor.RMQ)
	helper.FailOnError(err, "Fail on start RMQ connect")
	defer func() {
		log.Println("Close RMQ connect")
		err = r.CloseConnect()
		helper.FailOnError(err, "Fail on close RMQ connect")
	}()

	//sensor
	pzem004t, err := pzem.Setup(pzem.Config{Port: config.ServiceSensor.Device.Peacefair.Uart, Speed: 9600})
	helper.FailOnError(err, "Fail connect to pzem004t")

	ticker := time.NewTicker(time.Second * time.Duration(config.ServiceSensor.Device.Peacefair.Ticker))

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//send message to open RMQ connect by ticker
	go func() {
		ReadPzem004tAndSendRmq(pzem004t, r)
		for {
			select {
			case <-ticker.C:
				ReadPzem004tAndSendRmq(pzem004t, r)
			}
		}
	}()

	log.Printf("Got OS signal: %v. Exit...\n", <-osSignals)
}

func ReadPzem004tAndSendRmq(sensor pzem.Probe, r *rmq.RMQ) {
	jsonMessage, err := message.GetMessagePzem004tAsJson(sensor)
	helper.FailOnError(err, "Fail get message pzem004t")
	err = r.Send(jsonMessage)
	helper.FailOnError(err, "Fail on send RMQ message")
	log.Printf("Succes sentconfig RMQ message:%s", jsonMessage)
}
