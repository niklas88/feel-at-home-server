package clock

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Clock",
			Description: "Set device into clock mode"},
		ConfigFactory: deviceapi.EmptyConfigFactory,
		EffectFactory: deviceapi.WordClockEffectFactory(ClockEffectFactory)})
}

func ClockEffectFactory(l devices.WordClock) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		_, ok := config.(*deviceapi.EmptyConfig)
		if !ok {
			return errors.New("Not an empty Config")
		}

		return l.Clock()
	})
}
