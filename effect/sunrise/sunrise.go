package sunrise

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type SunriseConfig struct {
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Sunrise",
			Description: "A sunrise effect for color lamps"},
		ConfigFactory: func() effect.Config { return &SunriseConfig{"30ms"} },
		Effect:        effect.ColorLampEffect(SunriseEffect)})
}

func SunriseEffect(l lampbase.ColorLamp, conf effect.Config) error {
	sunriseConf, ok := conf.(*SunriseConfig)
	if !ok {
		return errors.New("Not a SunriseConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.Sunrise(delay)
}
