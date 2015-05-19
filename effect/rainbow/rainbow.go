package sunrise

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type RainbowConfig struct {
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Rainbow",
			Description: "A rainbow effect for color lamps"},
		ConfigFactory: func() effect.Config { return &RainbowConfig{"30ms"} },
		Effect:        effect.StripeLampEffect(RainbowEffect)})
}

func RainbowEffect(l lampbase.StripeLamp, conf effect.Config) error {
	sunriseConf, ok := conf.(*RainbowConfig)
	if !ok {
		return errors.New("Not a RainbowConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.Rainbow(delay)
}
