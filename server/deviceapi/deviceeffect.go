// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"github.com/niklas88/feel-at-home-server/devices"
)

type effectInfo struct {
	name          string
	description   string
	configFactory ConfigFactory
}

func (e *effectInfo) Name() string          { return e.name }
func (e *effectInfo) Description() string   { return e.description }
func (e *effectInfo) DefaultConfig() Config { return e.configFactory() }

type DeviceEffectFunc func(d devices.Device, config Config) error

type DeviceEffect struct {
	effectInfo
	applyToDevice DeviceEffectFunc
}

func (e *DeviceEffect) Apply(d devices.Device, config Config) error {
	return e.applyToDevice(d, config)
}

func (e *DeviceEffect) Compatible(d devices.Device) bool {
	return true
}

func NewDeviceEffect(name string, description string, f DeviceEffectFunc, cf ConfigFactory) *DeviceEffect {
	return &DeviceEffect{effectInfo{name, description, cf}, f}
}
