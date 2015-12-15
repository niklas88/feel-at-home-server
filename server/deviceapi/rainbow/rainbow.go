package rainbow

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Rainbow",
			Description: "A rainbow deviceapi for color lamps"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.StripeLampEffectFactory(RainbowEffectFactory)})
}

func RainbowEffectFactory(l devices.StripeLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		rainbowConf, ok := config.(*deviceapi.DelayConfig)
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
