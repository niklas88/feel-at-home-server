package strobe

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Stroboscope",
			Description: "Stroboscope"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.DimLampEffectFactory(StrobeEffectFactory)})
}

func StrobeEffectFactory(l lampbase.DimLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		strobeConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a StrobeConfig")
		}
		delay, err := time.ParseDuration(strobeConf.Delay)
		if err != nil {
			return err
		}
		return l.Stroboscope(delay)
	})
}
