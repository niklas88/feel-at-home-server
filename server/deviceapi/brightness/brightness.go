package brightness

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

type BrightnessConfig struct {
	Brightness uint8
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Brightness",
			Description: "Set brightness for your lamp"},
		ConfigFactory: func() deviceapi.Config {
			return &BrightnessConfig{255}
		},
		EffectFactory: deviceapi.DimLampEffectFactory(BrightnessEffectFactory)})
}

func BrightnessEffectFactory(l devices.DimLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		c, ok := config.(*BrightnessConfig)
		if !ok {
			return errors.New("Not a BrightnessConfig")
		}
		return l.Brightness(c.Brightness)
	})
}
