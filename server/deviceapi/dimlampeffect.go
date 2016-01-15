// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
)

type DimLampEffectFunc func(d devices.DimLamp, config Config) error

type DimLampEffect struct {
	effectInfo
	applyToDevice DimLampEffectFunc
}

func (e *DimLampEffect) Apply(dev devices.Device, config Config) error {
	d, ok := dev.(devices.DimLamp)
	if !ok {
		return errors.New("Incompatible device in DimLampEffect")
	}
	return e.applyToDevice(d, config)
}

func (e *DimLampEffect) Compatible(d devices.Device) bool {
	_, ok := d.(devices.DimLamp)
	return ok
}

func NewDimLampEffect(name string, description string, f DimLampEffectFunc, cf ConfigFactory) *DimLampEffect {
	return &DimLampEffect{effectInfo{name, description, cf}, f}
}
