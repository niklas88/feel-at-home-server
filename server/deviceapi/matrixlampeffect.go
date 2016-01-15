// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
)

type MatrixLampEffectFunc func(d devices.MatrixLamp, config Config) error

type MatrixLampEffect struct {
	effectInfo
	applyToDevice MatrixLampEffectFunc
}

func (e *MatrixLampEffect) Apply(dev devices.Device, config Config) error {
	d, ok := dev.(devices.MatrixLamp)
	if !ok {
		return errors.New("Incompatible device in MatrixLampEffect")
	}
	return e.applyToDevice(d, config)
}

func (e *MatrixLampEffect) Compatible(d devices.Device) bool {
	_, ok := d.(devices.MatrixLamp)
	return ok
}

func NewMatrixLampEffect(name string, description string, f MatrixLampEffectFunc, cf ConfigFactory) *MatrixLampEffect {
	return &MatrixLampEffect{effectInfo{name, description, cf}, f}
}
