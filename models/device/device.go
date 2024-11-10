package device

import (
	"fmt"
	"strings"

	devicepart "github.com/poriamsz55/BoosterPump-webapp/models/device_part"
)

type Converter int

const (
	WithoutConverter Converter = iota
	OneToOne
	TwoToTwo
)

var converterNames = []string{
	"بدون تبدیل",
	"تبدیل در دهش",
	"تبدیل دو طرفه",
}

func (c Converter) String() string {
	if c < WithoutConverter || c > TwoToTwo {
		return "Unknown"
	}
	return converterNames[c]
}

func ConverterFromValue(value int) (Converter, error) {
	if value < int(WithoutConverter) || value > int(TwoToTwo) {
		return 0, fmt.Errorf("invalid value: %d", value)
	}
	return Converter(value), nil
}

func ConverterFromName(name string) (Converter, error) {
	for i, n := range converterNames {
		if strings.EqualFold(n, name) {
			return Converter(i), nil
		}
	}
	return 0, fmt.Errorf("invalid name: %s", name)
}

type Device struct {
	Id             int                      `json:"id"`
	Name           string                   `json:"string"`
	Converter      Converter                `json:"converter"`
	Filter         bool                     `json:"filter"`
	DevicePartList []*devicepart.DevicePart `json:"device_part"`
	Price          uint64                   `json:"price"`
}

func NewEmptyDevice() *Device {
	return &Device{
		DevicePartList: []*devicepart.DevicePart{},
	}
}

func NewDevice(name string, converter Converter, filter bool) *Device {
	return &Device{
		Name:           name,
		Converter:      converter,
		Filter:         filter,
		DevicePartList: []*devicepart.DevicePart{},
	}
}

func (d *Device) UpdatePrice() {
	d.Price = 0
	for _, dp := range d.DevicePartList {
		d.Price += dp.Price
	}
}
