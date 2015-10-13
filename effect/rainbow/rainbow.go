package rainbow

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Rainbow",
			Description: "A rainbow effect for color lamps"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.StripeLampEffectFactory(RainbowEffectFactory)})
}

func RainbowEffectFactory(l lampbase.StripeLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		sunriseConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a RainbowConfig")
		}

		delay, err := time.ParseDuration(sunriseConf.Delay)
		if err != nil {
			return err
		}
		return l.Rainbow(delay)
	})
}
