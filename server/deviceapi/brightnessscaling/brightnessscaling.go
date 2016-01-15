package brightness

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

type BrightnessConfig struct {
	Brightness uint8
}

func applyToDevice(l devices.DimLamp, config deviceapi.Config) error {
	c, ok := config.(*BrightnessConfig)
	if !ok {
		return errors.New("Not a BrightnessConfig")
	}
	return l.BrightnessScaling(c.Brightness)
}

func init() {
	deviceapi.DefaultRegistry.Register(
		deviceapi.NewDimLampEffect(
			"Brightnessscaling",
			"Set brightness scaling for your lamp",
			applyToDevice,
			func() deviceapi.Config {
				return &BrightnessConfig{255}
			}))
}
