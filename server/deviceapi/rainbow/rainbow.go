package rainbow

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewStripeLampEffect(
		"Rainbow",
		"A rainbow effect for stripe lamps",
		applyToDevice,
		func() deviceapi.Config { return &deviceapi.DelayConfig{"10ms"} }))
}

func applyToDevice(l devices.StripeLamp, config deviceapi.Config) error {
	rainbowConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a RainbowConfig")
	}

	delay, err := time.ParseDuration(rainbowConf.Delay)
	if err != nil {
		return err
	}
	return l.Rainbow(delay)
}
