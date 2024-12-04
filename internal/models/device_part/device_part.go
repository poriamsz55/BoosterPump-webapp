package devicepart

import (
	"math"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

type DevicePart struct {
	Id       int        `json:"id"`
	DeviceId int        `json:"device_id"`
	Price    uint64     `json:"price"`
	Count    float64    `json:"count"`
	Part     *part.Part `json:"part"`
}

func NewEmptyDevicePart() *DevicePart {
	return &DevicePart{}
}

func NewDevicePart(deviceId int, count float64, prt *part.Part) *DevicePart {
	return &DevicePart{
		DeviceId: deviceId,
		Count:    count,
		Part:     prt,
	}
}

func (dp *DevicePart) UpdatePrice() {
	dp.Price = uint64(math.Ceil(float64(dp.Count) * float64(dp.Part.Price)))
}
