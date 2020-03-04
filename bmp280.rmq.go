package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dark705/go_sensor/internal/message"
	"github.com/dark705/go_sensor/internal/rmq"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceSensor struct {
		RMQ    rmq.Config `yaml:"rabbitmq"`
		Device struct {
			BMP280 struct {
				Ticker int `yaml:"ticker"`
			} `yaml:"bmp280"`
		} `yaml:"device"`
	} `yaml:"service_sensor"`
}

func main() {
	//config
	var config Config
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	failOnError(err, "Cant read config file")
	err = yaml.Unmarshal(yamlFile, &config)
	failOnError(err, "Cant read format of config file")

	//create RMQ connect
	r, err := rmq.NewRMQ(config.ServiceSensor.RMQ)
	failOnError(err, "Fail on start RMQ connect")
	defer func() {
		log.Println("Close RMQ connect")
		err = r.CloseConnect()
		failOnError(err, "Fail on close RMQ connect")
	}()

	//connect to I2C bus
	i2cBus, err := i2c.NewI2C(0x76, 1)
	_ = logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	failOnError(err, "Failed on i2c bus")
	defer func() {
		log.Println("Close I2C connect")
		err := i2cBus.Close()
		failOnError(err, "Fail on close I2C connect")
	}()

	//sensor
	bmp280, err := bsbmp.NewBMP(bsbmp.BMP280, i2cBus)
	failOnError(err, "Failed on bmp280")
	_ = logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)

	ticker := time.NewTicker(time.Second * time.Duration(config.ServiceSensor.Device.BMP280.Ticker))

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//send message to open RMQ connect by ticker
	go func() {
		ReadAndSend(r, bmp280)
		for {
			select {
			case <-ticker.C:
				ReadAndSend(r, bmp280)
			}
		}
	}()

	log.Printf("Got OS signal: %v. Exit...\n", <-osSignals)
}

func ReadAndSend(r *rmq.RMQ, bmp280 *bsbmp.BMP) {
	jsonMessage, err := message.GetMessageAsJson(bmp280, message.SensorBMP280)
	failOnError(err, "Fail get message bmp280")
	err = r.Send(jsonMessage)
	failOnError(err, "Fail on send RMQ message")
	log.Printf("Succes sentconfig RMQ message:%s", jsonMessage)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
