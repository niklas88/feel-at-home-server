package whitefade

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(deviceapi.NewDimLampEffect(
		"Whitefade",
		"White fading effect",
		applyToDevice,
		func() deviceapi.Config { return &deviceapi.DelayConfig{"10ms"} }))
}

func applyToDevice(l devices.DimLamp, config deviceapi.Config) error {
	strobeConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a WhiteFadeConfig")
	}
	delay, err := time.ParseDuration(strobeConf.Delay)
	if err != nil {
		return err
	}
	return l.Fade(delay, 255)
}
