package strobe

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type StrobeConfig struct {
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Stroboscope",
			Description: "Stroboscope"},
		ConfigFactory: func() effect.Config {
			return &StrobeConfig{"30ms"}
		},
		Effect: effect.DimLampEffect(StrobeEffect)})
}

func StrobeEffect(l lampbase.DimLamp, conf effect.Config) error {
	strobeConf, ok := conf.(*StrobeConfig)
	if !ok {
		return errors.New("Not a StrobeConfig")
	}
	delay, err := time.ParseDuration(strobeConf.Delay)
	if err != nil {
		return err
	}
	return l.Stroboscope(delay)
}
