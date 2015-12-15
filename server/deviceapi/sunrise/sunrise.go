package sunrise

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

type SunriseConfig struct {
	Delay string
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Sunrise",
			Description: "A sunrise deviceapi for color lamps"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.ColorLampEffectFactory(SunriseEffectFactory)})
}

func SunriseEffectFactory(l devices.ColorLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		sunriseConf, ok := config.(*deviceapi.DelayConfig)
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
