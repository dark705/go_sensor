package message

import (
	"encoding/json"
	"fmt"
	"time"

	nut "github.com/robbiet480/go.nut"
)

const (
	SensorUPS string = "ups"
)

type UPSMessage struct {
	Sensor string           `json:"sensor"`
	Date   time.Time        `json:"date"`
	Data   []UPSMessageData `json:"data"`
}

type UPSMessageData struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

func GetMessageUPS(nutClient nut.Client) (UPSMessage, error) {
	upsMessage := UPSMessage{Sensor: SensorUPS, Date: time.Now()}

	upsList, err := nutClient.GetUPSList()
	if err != nil {
		return upsMessage, err
	}

	UPSMessageDataList := make([]UPSMessageData, 0, len(upsList))

	for _, ups := range upsList {
		upsInfo := UPSMessageData{Name: ups.Name, Data: make(map[string]string)}
		for _, variable := range ups.Variables {
			upsInfo.Data[variable.Name] = fmt.Sprintf("%v", variable.Value)
		}
		UPSMessageDataList = append(UPSMessageDataList, upsInfo)
	}
	upsMessage.Data = UPSMessageDataList

	return upsMessage, nil
}

func GetMessageUPSAsJSON(nutClient nut.Client) ([]byte, error) {
	message, err := GetMessageUPS(nutClient)
	jsonMessage, _ := json.Marshal(message)
	return jsonMessage, err
}
