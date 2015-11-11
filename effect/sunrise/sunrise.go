package sunrise

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/effect"
	"github.com/niklas88/feel-at-home-server/lampbase"
	"time"
)

type SunriseConfig struct {
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Sunrise",
			Description: "A sunrise effect for color lamps"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.ColorLampEffectFactory(SunriseEffectFactory)})
}

func SunriseEffectFactory(l lampbase.ColorLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		sunriseConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a SunriseConfig")
		}

		delay, err := time.ParseDuration(sunriseConf.Delay)
		if err != nil {
			return err
		}
		return l.Sunrise(delay)
	})
}
