package strobe

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Stroboscope",
			Description: "Stroboscope"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.DimLampEffectFactory(StrobeEffectFactory)})
}

func StrobeEffectFactory(l devices.DimLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		strobeConf, ok := config.(*deviceapi.DelayConfig)
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
