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

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
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

	//connect to I2C bus
	i2cBus, err := i2c.NewI2C(0x76, 1)
	_ = logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	helper.FailOnError(err, "Failed on i2c bus")
	defer func() {
		log.Println("Close I2C connect")
		err := i2cBus.Close()
		helper.FailOnError(err, "Fail on close I2C connect")
	}()

	//sensor
	bmp280, err := bsbmp.NewBMP(bsbmp.BMP280, i2cBus)
	helper.FailOnError(err, "Failed on bmp280")
	_ = logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)

	ticker := time.NewTicker(time.Second * time.Duration(config.ServiceSensor.Device.BMP280.Ticker))

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//send message to open RMQ connect by ticker
	go func() {
		ReadBmp280AndSendRmq(bmp280, r)
		for {
			select {
			case <-ticker.C:
				ReadBmp280AndSendRmq(bmp280, r)
			}
		}
	}()

	log.Printf("Got OS signal: %v. Exit...\n", <-osSignals)
	ticker.Stop()
}

func ReadBmp280AndSendRmq(bmp280 *bsbmp.BMP, r *rmq.RMQ) {
	jsonMessage, err := message.GetMessageBmp280AsJson(bmp280, message.SensorBMP280)
	helper.FailOnError(err, "Fail get message bmp280")
	err = r.Send(jsonMessage)
	helper.FailOnError(err, "Fail on send RMQ message")
	log.Printf("Succes sent RMQ message:%s", jsonMessage)
}
