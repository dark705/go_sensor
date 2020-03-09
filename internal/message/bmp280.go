package message

import (
	"encoding/json"
	"time"

	"github.com/d2r2/go-bsbmp"
)

const (
	SensorBMP280 string = "bmp280"
)

type BmpMessage struct {
	Sensor string         `json:"sensor"`
	Date   time.Time      `json:"date"`
	Data   BmpMessageData `json:"data"`
}

type BmpMessageData struct {
	Pressure    float32 `json:"pressure"`
	Temperature float32 `json:"temperature"`
}

func GetMessageBmp280(sensor *bsbmp.BMP, name string) (BmpMessage, error) {
	// Read atmospheric pressure in mmHg
	pressure, err := sensor.ReadPressureMmHg(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return BmpMessage{}, err
	}

	// Read temperature in celsius degree
	temperature, err := sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return BmpMessage{}, err
	}

	return BmpMessage{
		Sensor: name,
		Date:   time.Now(),
		Data:   BmpMessageData{Pressure: pressure, Temperature: temperature},
	}, nil
}

func GetMessageBmp280AsJson(sensor *bsbmp.BMP, name string) ([]byte, error) {
	message, err := GetMessageBmp280(sensor, name)
	jsonMessage, _ := json.Marshal(message)
	return jsonMessage, err
}
