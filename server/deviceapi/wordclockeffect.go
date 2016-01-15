// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
)

type WordClockEffectFunc func(d devices.WordClock, config Config) error

type WordClockEffect struct {
	effectInfo
	applyToDevice WordClockEffectFunc
}

func (e *WordClockEffect) Apply(dev devices.Device, config Config) error {
	d, ok := dev.(devices.WordClock)
	if !ok {
		return errors.New("Incompatible device in WordClockEffect")
	}
	return e.applyToDevice(d, config)
}

func (e *WordClockEffect) Compatible(d devices.Device) bool {
	_, ok := d.(devices.WordClock)
	return ok
}

func NewWordClockEffect(name string, description string, f WordClockEffectFunc, cf ConfigFactory) *WordClockEffect {
	return &WordClockEffect{effectInfo{name, description, cf}, f}
}
