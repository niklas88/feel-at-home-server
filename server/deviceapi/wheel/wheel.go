package wheel

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewColorLampEffect(
		"Wheel",
		"A color changing effect for color lamps",
		applyToDevice,
		func() deviceapi.Config { return &deviceapi.DelayConfig{"10ms"} }))
}

func applyToDevice(l devices.ColorLamp, config deviceapi.Config) error {
	wheelConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a WheelConf")
	}

	delay, err := time.ParseDuration(wheelConf.Delay)
	if err != nil {
		return err
	}
	return l.ColorWheel(delay)
}
