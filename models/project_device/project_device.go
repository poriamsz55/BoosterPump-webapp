package projectdevice

import (
	"math"

	"github.com/poriamsz55/BoosterPump-webapp/models/device"
)

type ProjectDevice struct {
	Id        int            `json:"id"`
	ProjectId int            `json:"project_id"`
	Count     float32        `json:"count"`
	Device    *device.Device `json:"device"`
	Price     uint64         `json:"price"`
}

func NewEmptyProjectDevice() *ProjectDevice {
	return &ProjectDevice{}
}

func NewProjectDevice(prjId int, count float32, device *device.Device) *ProjectDevice {
	return &ProjectDevice{
		ProjectId: prjId,
		Count:     count,
		Device:    device,
	}
}

func (pd *ProjectDevice) UpdatePrice() {
	pd.Price = uint64(math.Ceil(float64(pd.Count * float32(pd.Device.Price))))
}
