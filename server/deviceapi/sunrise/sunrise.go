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
	deviceapi.DefaultRegistry.Register(deviceapi.NewColorLampEffect(
		"Sunrise",
		"A sunrise effect for color lamps",
		applyToDevice,
		func() deviceapi.Config { return deviceapi.DelayConfig{"10ms"} }))
}

func applyToDevice(l devices.ColorLamp, config deviceapi.Config) error {
	sunriseConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a SunriseConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.Sunrise(delay)
}
