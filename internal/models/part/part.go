package part

type Part struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Price    uint64 `json:"price"`
	Size     string `json:"size"`
	Material string `json:"material"`
	Brand    string `json:"brand"`
}

func NewEmptyPart() *Part {
	return &Part{}
}

func NewPart(name, size, material, brand string, price uint64) *Part {
	return &Part{
		Name:     name,
		Size:     size,
		Material: material,
		Brand:    brand,
		Price:    price,
	}
}

type PartJson struct {
	Id    string `json:"id"`
	Count string `json:"count"`
}

type PartReq struct {
	Id    int     `json:"id"`
	Count float64 `json:"count"`
}
