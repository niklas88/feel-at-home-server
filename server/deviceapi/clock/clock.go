package clock

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewWordClockEffect(
		"Clock",
		"Set device into clock mode",
		ApplyToDevice,
		func() deviceapi.Config { return &deviceapi.EmptyConfig{} }))
}

func ApplyToDevice(l devices.WordClock, config deviceapi.Config) error {
	_, ok := config.(*deviceapi.EmptyConfig)
	if !ok {
		return errors.New("Not an empty Config")
	}

	return l.Clock()
}
