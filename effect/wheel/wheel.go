package sunrise

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type WheelConfig struct {
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Wheel",
			Description: "A color wheel effect for color lamps"},
		ConfigFactory: func() effect.Config { return &WheelConfig{"30ms"} },
		Effect:        effect.ColorLampEffect(WheelEffect)})
}

func WheelEffect(l lampbase.ColorLamp, conf effect.Config) error {
	sunriseConf, ok := conf.(*WheelConfig)
	if !ok {
		return errors.New("Not a WheelConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.ColorWheel(delay)
}
