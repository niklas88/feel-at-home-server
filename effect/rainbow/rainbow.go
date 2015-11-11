package rainbow

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/effect"
	"github.com/niklas88/feel-at-home-server/lampbase"
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
		rainbowConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a RainbowConfig")
		}

		delay, err := time.ParseDuration(rainbowConf.Delay)
		if err != nil {
			return err
		}
		return l.Rainbow(delay)
	})
}
