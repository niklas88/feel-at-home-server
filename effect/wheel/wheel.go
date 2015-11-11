package wheel

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/effect"
	"github.com/niklas88/feel-at-home-server/lampbase"
	"time"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Wheel",
			Description: "A color chaning effect for color lamps"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.ColorLampEffectFactory(WheelEffectFactory)})
}

func WheelEffectFactory(l lampbase.ColorLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		wheelConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a WheelConf")
		}

		delay, err := time.ParseDuration(wheelConf.Delay)
		if err != nil {
			return err
		}
		return l.ColorWheel(delay)
	})
}
