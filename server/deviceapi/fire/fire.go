package fire

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

type FireConfig struct {
	Delay string
	Cooling uint8
	Spark uint8
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Fire",
			Description: "A fire deviceapi for color lamps"},
		ConfigFactory: func() deviceapi.Config { return &FireConfig{"32ms", 71, 80} },
		EffectFactory: deviceapi.StripeLampEffectFactory(FireEffectFactory)})
}

func FireEffectFactory(l devices.StripeLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		fireConf, ok := config.(*FireConfig)
		if !ok {
			return errors.New("Not a FireConfig")
		}

		delay, err := time.ParseDuration(fireConf.Delay)
		if err != nil {
			return err
		}
		return l.Fire(delay, fireConf.Cooling, fireConf.Spark)
	})
}
