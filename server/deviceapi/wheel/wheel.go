package wheel

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Wheel",
			Description: "A color chaning deviceapi for color lamps"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.ColorLampEffectFactory(WheelEffectFactory)})
}

func WheelEffectFactory(l devices.ColorLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		wheelConf, ok := config.(*deviceapi.DelayConfig)
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
