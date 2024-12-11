package project

import (
	"time"

	extraprice "github.com/poriamsz55/BoosterPump-webapp/internal/models/extra_price"
	projectdevice "github.com/poriamsz55/BoosterPump-webapp/internal/models/project_device"
)

type Project struct {
	Id                int                            `json:"id"`
	Name              string                         `json:"name"`
	Price             uint64                         `json:"price"`
	ProjectDeviceList []*projectdevice.ProjectDevice `json:"project_device"`
	ExtraPriceList    []*extraprice.ExtraPrice       `json:"extra_price"`
	ModifiedAt        time.Time                      `json:"modified_at"`
}

func NewEmptyProject() *Project {
	return &Project{
		ProjectDeviceList: []*projectdevice.ProjectDevice{},
		ExtraPriceList:    []*extraprice.ExtraPrice{},
	}
}

func NewProject(name string) *Project {
	return &Project{
		Name:              name,
		ProjectDeviceList: []*projectdevice.ProjectDevice{},
		ExtraPriceList:    []*extraprice.ExtraPrice{},
	}
}

func (p *Project) UpdatePrice() {

	p.Price = 0
	if p.ProjectDeviceList != nil {
		for _, pd := range p.ProjectDeviceList {
			pd.UpdatePrice()
			p.Price += pd.Price
		}
	}

	if p.ExtraPriceList != nil {
		for _, ep := range p.ExtraPriceList {
			p.Price += ep.Price
		}
	}

}
