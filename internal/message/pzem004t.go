package message

import (
	"encoding/json"
	"time"

	"github.com/dark705/go_sensor/internal/helper"
	"github.com/dark705/pzem-004t-v3/pzem"
)

const SensorPzem004t string = "pzem004t.v3"

type PzemMessage struct {
	Sensor string          `json:"sensor"`
	Date   time.Time       `json:"date"`
	Data   PzemMessageData `json:"data"`
}

type PzemMessageData struct {
	Voltage     float32 `json:"voltage"`
	Intensity   float32 `json:"current"`
	Power       float32 `json:"active"`
	Frequency   float32 `json:"frequency"`
	Energy      float32 `json:"energy"`
	PowerFactor float32 `json:"powerFactor"`
}

func GetMessagePzem004t(sensor pzem.Probe) (PzemMessage, error) {
	voltage, err := sensor.Voltage()
	helper.FailOnError(err, "Fail on read voltage")

	intensity, err := sensor.Intensity()
	helper.FailOnError(err, "Fail on read intensity")

	power, err := sensor.Power()
	helper.FailOnError(err, "Fail on read power")

	frequency, err := sensor.Frequency()
	helper.FailOnError(err, "Fail on read frequency")

	energy, err := sensor.Energy()
	helper.FailOnError(err, "Fail on read energy")

	powerFactor, err := sensor.PowerFactor()
	helper.FailOnError(err, "Fail on read powerFactor")

	return PzemMessage{
		Sensor: SensorPzem004t,
		Date:   time.Now(),
		Data: PzemMessageData{
			Voltage:     voltage,
			Intensity:   intensity,
			Power:       power,
			Frequency:   frequency,
			Energy:      energy,
			PowerFactor: powerFactor,
		},
	}, nil
}

func GetMessagePzem004tAsJson(sensor pzem.Probe) ([]byte, error) {
	message, err := GetMessagePzem004t(sensor)
	jsonMessage, _ := json.Marshal(message)
	return jsonMessage, err
}
