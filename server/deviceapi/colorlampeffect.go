// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
)

type ColorLampEffectFunc func(d devices.ColorLamp, config Config) error

type ColorLampEffect struct {
	effectInfo
	applyToDevice ColorLampEffectFunc
}

func (e *ColorLampEffect) Apply(dev devices.Device, config Config) error {
	d, ok := dev.(devices.ColorLamp)
	if !ok {
		return errors.New("Incompatible device in ColorLampEffect")
	}
	return e.applyToDevice(d, config)
}

func (e *ColorLampEffect) Compatible(d devices.Device) bool {
	_, ok := d.(devices.ColorLamp)
	return ok
}

func NewColorLampEffect(name string, description string, f ColorLampEffectFunc, cf ConfigFactory) *ColorLampEffect {
	return &ColorLampEffect{effectInfo{name, description, cf}, f}
}
