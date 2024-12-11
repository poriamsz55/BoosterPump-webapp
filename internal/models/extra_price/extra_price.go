package extraprice

import "time"

type ExtraPrice struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	ProjectId  int       `json:"project_id"`
	Price      uint64    `json:"price"`
	ModifiedAt time.Time `json:"modified_at"`
}

func NewEmptyExtraPrice() *ExtraPrice {
	return &ExtraPrice{}
}

func NewExtraPrice(projectId int, name string, price uint64) *ExtraPrice {
	return &ExtraPrice{
		ProjectId: projectId,
		Name:      name,
		Price:     price,
	}
}
