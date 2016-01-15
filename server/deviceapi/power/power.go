package power

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

type PowerConfig struct {
	Power bool
}

func applyToDevice(d devices.Device, config deviceapi.Config) error {
	power, ok := config.(*PowerConfig)
	if !ok {
		return errors.New("Not a PowerConfig")
	}

	return d.Power(power.Power)
}

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewDeviceEffect("Power", "Turns devices on or off", applyToDevice, func() deviceapi.Config { return &PowerConfig{false} }))
}
