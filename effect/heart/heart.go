package heart

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Heart",
			Description: "Set device into heart mode"},
		ConfigFactory: effect.EmptyConfigFactory,
		EffectFactory: effect.MatrixLampEffectFactory(HeartEffectFactory)})
}

func HeartEffectFactory(l lampbase.MatrixLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		_, ok := config.(*effect.EmptyConfig)
		if !ok {
			return errors.New("Not an empty Config")
		}

		return l.Heart()
	})
}
