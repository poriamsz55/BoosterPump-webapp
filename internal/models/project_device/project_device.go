package projectd

import (
	"math"
	"time"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
)

type ProjectDevice struct {
	Id         int            `json:"id"`
	ProjectId  int            `json:"project_id"`
	Count      float64        `json:"count"`
	Device     *device.Device `json:"device"`
	Price      uint64         `json:"price"`
	ModifiedAt time.Time      `json:"modified_at"`
}

func NewEmptyProjectDevice() *ProjectDevice {
	return &ProjectDevice{}
}

func NewProjectDevice(prjId int, count float64, d *device.Device) *ProjectDevice {
	return &ProjectDevice{
		ProjectId: prjId,
		Count:     count,
		Device:    d,
	}
}

func (pd *ProjectDevice) UpdatePrice() {
	pd.Price = uint64(math.Ceil(float64(pd.Count * float64(pd.Device.Price))))
}
