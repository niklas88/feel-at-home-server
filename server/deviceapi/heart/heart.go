package heart

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewMatrixLampEffect(
		"Heart",
		"Set device into heart mode",
		applyToDevice,
		func() deviceapi.Config { return &deviceapi.EmptyConfig{} }))
}

func applyToDevice(l devices.MatrixLamp, config deviceapi.Config) error {
	_, ok := config.(*deviceapi.EmptyConfig)
	if !ok {
		return errors.New("Not an empty Config")
	}

	return l.Heart()
}
