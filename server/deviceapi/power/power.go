package power

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

type PowerConfig struct {
	Power bool
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Power",
			Description: "Turn your device on and off"},
		ConfigFactory: func() deviceapi.Config {
			return &PowerConfig{true}
		},
		EffectFactory: deviceapi.DeviceEffectFactory(PowerEffect)})
}

func PowerEffect(l devices.Device) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		power, ok := config.(*PowerConfig)
		if !ok {
			return errors.New("Not a PowerConfig")
		}

		return l.Power(power.Power)
	})
}
