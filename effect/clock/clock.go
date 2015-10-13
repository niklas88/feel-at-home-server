package clock

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Clock",
			Description: "Set device into clock mode"},
		ConfigFactory: effect.EmptyConfigFactory,
		EffectFactory: effect.WordClockEffectFactory(ClockEffectFactory)})
}

func ClockEffectFactory(l lampbase.WordClock) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		_, ok := config.(*effect.EmptyConfig)
		if !ok {
			return errors.New("Not an empty Config")
		}

		return l.Clock()
	})
}
