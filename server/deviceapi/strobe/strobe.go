package strobe

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewDimLampEffect(
		"Stroboscope",
		"A stroboscope effect",
		applyToDevice,
		func() deviceapi.Config { return &deviceapi.DelayConfig{"40ms"} }))
}

func applyToDevice(l devices.DimLamp, config deviceapi.Config) error {
	strobeConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a StrobeConfig")
	}
	delay, err := time.ParseDuration(strobeConf.Delay)
	if err != nil {
		return err
	}
	return l.Stroboscope(delay)
}
