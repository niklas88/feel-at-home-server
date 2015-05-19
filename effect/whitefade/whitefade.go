package strobe

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type WhiteFadeConfig struct {
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "White Fade",
			Description: "White fading effect"},
		ConfigFactory: func() effect.Config {
			return &WhiteFadeConfig{"30ms"}
		},
		Effect: effect.DimLampEffect(WhiteFadeEffect)})
}

func WhiteFadeEffect(l lampbase.DimLamp, conf effect.Config) error {
	strobeConf, ok := conf.(*WhiteFadeConfig)
	if !ok {
		return errors.New("Not a WhiteFadeConfig")
	}
	delay, err := time.ParseDuration(strobeConf.Delay)
	if err != nil {
		return err
	}
	return l.Fade(delay, 255)
}
