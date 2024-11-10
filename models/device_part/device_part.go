package devicepart

import (
	"math"

	"github.com/poriamsz55/BoosterPump-webapp/models/part"
)

type DevicePart struct {
	Id       int        `json:"id"`
	DeviceId int        `json:"device_id"`
	Price    uint64     `json:"price"`
	Count    float32    `json:"count"`
	Part     *part.Part `json:"aprt"`
}

func NewEmptyDevicePart() *DevicePart {
	return &DevicePart{}
}

func NewDevicePart(deviceId int, count float32, prt *part.Part) *DevicePart {
	return &DevicePart{
		DeviceId: deviceId,
		Count:    count,
		Part:     prt,
	}
}

func (dp *DevicePart) UpdatePrice() {
	dp.Price = uint64(math.Ceil(float64(dp.Count * float32(dp.Part.Price))))
}
