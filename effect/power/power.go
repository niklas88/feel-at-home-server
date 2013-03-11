package power

import (
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type PowerConfig struct {
	Power bool
}

type PowerEffect struct {
	power PowerConfig
	lamp  lampbase.Device
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Power",
			Description: "Turn your device on and off"},
		ConfigFactory: func() effect.Config {
			return &PowerConfig{true}
		},
		Factory: effect.DeviceEffectFactory(NewPowerEffect)})
}

func NewPowerEffect(l lampbase.Device) effect.Effect {
	return &PowerEffect{PowerConfig{true}, l}
}

func (p *PowerEffect) Configure(c effect.Config) {
	p.power = *c.(*PowerConfig)
}

func (p *PowerEffect) Apply() (time.Duration, error) {
	err := p.lamp.Power(p.power.Power)
	return -1, err
}
