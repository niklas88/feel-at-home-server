package power

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/effect"
	"github.com/niklas88/feel-at-home-server/lampbase"
)

type PowerConfig struct {
	Power bool
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Power",
			Description: "Turn your device on and off"},
		ConfigFactory: func() effect.Config {
			return &PowerConfig{true}
		},
		EffectFactory: effect.DeviceEffectFactory(PowerEffect)})
}

func PowerEffect(l lampbase.Device) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		power, ok := config.(*PowerConfig)
		if !ok {
			return errors.New("Not a PowerConfig")
		}

		return l.Power(power.Power)
	})
}
