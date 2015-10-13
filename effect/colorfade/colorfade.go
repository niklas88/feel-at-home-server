package colorefade

import (
	"errors"
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type ColorfadeConfig struct {
	Color hexcolor.Hex
	Delay string
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Colorfade",
			Description: "Fades with Color"},
		ConfigFactory: func() effect.Config { return &ColorfadeConfig{"#ffffff", "15ms"} },
		EffectFactory: effect.ColorLampEffectFactory(ColorFadeEffect)})
}

func ColorFadeEffect(l lampbase.ColorLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		colorfadeConf, ok := config.(*ColorfadeConfig)
		if !ok {
			return errors.New("Not a ColorFadeConfig")
		}

		delay, err := time.ParseDuration(colorfadeConf.Delay)
		if err != nil {
			return err
		}

		m := color.RGBAModel
		return l.ColorFade(delay, m.Convert(colorfadeConf.Color).(color.RGBA))
	})
}
