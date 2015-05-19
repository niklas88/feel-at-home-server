package brightness

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
)

type BrightnessConfig struct {
	Brightness uint8
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Brightness",
			Description: "Set brightness for your lamp"},
		ConfigFactory: func() effect.Config {
			return &BrightnessConfig{255}
		},
		Effect: effect.DimLampEffect(BrightnessEffect)})
}

func BrightnessEffect(l lampbase.DimLamp, config effect.Config) error {
	c, ok := config.(*BrightnessConfig)
	if !ok {
		return errors.New("Not a BrightnessConfig")
	}
	return l.SetBrightness(c.Brightness)
}
