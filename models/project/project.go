package project

import (
	extraprice "github.com/poriamsz55/BoosterPump-webapp/models/extra_price"
	projectdevice "github.com/poriamsz55/BoosterPump-webapp/models/project_device"
)

type Project struct {
	Id                int                            `json:"id"`
	Name              string                         `json:"name"`
	Price             uint64                         `json:"price"`
	ProjectDeviceList []*projectdevice.ProjectDevice `json:"project_device"`
	ExtraPriceLsit    []*extraprice.ExtraPrice       `json:"extra_price"`
}

func NewEmptyProject() *Project {
	return &Project{
		ProjectDeviceList: []*projectdevice.ProjectDevice{},
		ExtraPriceLsit:    []*extraprice.ExtraPrice{},
	}
}

func NewProject(name string) *Project {
	return &Project{
		Name:              name,
		ProjectDeviceList: []*projectdevice.ProjectDevice{},
		ExtraPriceLsit:    []*extraprice.ExtraPrice{},
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

	if p.ExtraPriceLsit != nil {
		for _, ep := range p.ExtraPriceLsit {
			p.Price += ep.Price
		}
	}

}
