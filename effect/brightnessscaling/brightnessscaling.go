package brightnessscaling

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
			Name:        "Brightness Scaling",
			Description: "Set brightness scaling for your lamp"},
		ConfigFactory: func() effect.Config {
			return &BrightnessConfig{255}
		},
		EffectFactory: effect.DimLampEffectFactory(BrightnessScalingEffectFactory)})
}

func BrightnessScalingEffectFactory(l lampbase.DimLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		c, ok := config.(*BrightnessConfig)
		if !ok {
			return errors.New("Not a BrightnessConfig")
		}
		return l.BrightnessScaling(c.Brightness)
	})
}
