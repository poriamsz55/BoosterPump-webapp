package part

type Part struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Price    uint64 `json:"price"`
	Size     string `json:"size"`
	Material string `json:"matrial"`
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
