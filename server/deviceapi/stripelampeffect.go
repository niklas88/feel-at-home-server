// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
)

type StripeLampEffectFunc func(d devices.StripeLamp, config Config) error

type StripeLampEffect struct {
	effectInfo
	applyToDevice StripeLampEffectFunc
}

func (e *StripeLampEffect) Apply(dev devices.Device, config Config) error {
	d, ok := dev.(devices.StripeLamp)
	if !ok {
		return errors.New("Incompatible device in StripeLampEffect")
	}
	return e.applyToDevice(d, config)
}

func (e *StripeLampEffect) Compatible(d devices.Device) bool {
	_, ok := d.(devices.StripeLamp)
	return ok
}

func NewStripeLampEffect(name string, description string, f StripeLampEffectFunc, cf ConfigFactory) *StripeLampEffect {
	return &StripeLampEffect{effectInfo{name, description, cf}, f}
}
