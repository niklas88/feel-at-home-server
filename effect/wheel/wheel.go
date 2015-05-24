package wheel

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Wheel",
			Description: "A color wheel effect for color lamps"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.ColorLampEffectFactory(NewWheelEffect)})
}

func NewWheelEffect(l lampbase.ColorLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		sunriseConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a WheelConfig")
		}

		delay, err := time.ParseDuration(sunriseConf.Delay)
		if err != nil {
			return err
		}
		return l.ColorWheel(delay)
	})
}
