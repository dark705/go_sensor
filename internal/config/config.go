package сonfig

import (
	"io/ioutil"

	"github.com/dark705/go_sensor/internal/helper"
	"github.com/dark705/go_sensor/internal/rmq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceSensor struct {
		RMQ    rmq.Config `yaml:"rabbitmq"`
		Device struct {
			BMP280 struct {
				Ticker int `yaml:"ticker"`
			} `yaml:"bmp280"`
			Peacefair struct {
				Uart    string `yaml:"uart"`
				Ticker  int    `yaml:"ticker"`
				Timeout int    `yaml:"timeout"`
			} `yaml:"peacefair"`
		} `yaml:"device"`
	} `yaml:"service_sensor"`
}

func ReadFromFile(file string) (conf Config) {
	if file == "" {
		file = "config/config.yaml" //default file location
	}
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	helper.FailOnError(err, "Cant read config file")
	err = yaml.Unmarshal(yamlFile, &conf)
	helper.FailOnError(err, "Cant read format of config file")
	return conf
}
