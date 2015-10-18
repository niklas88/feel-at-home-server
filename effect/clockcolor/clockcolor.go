package clockcolor

import (
	"errors"
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
)

type ClockColorConfig struct {
	Color hexcolor.Hex
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Clock Color",
			Description: "Set the color for your clock after clock effect selected"},
		ConfigFactory: func() effect.Config {
			return &ClockColorConfig{"#ffffff"}
		},
		EffectFactory: effect.WordClockEffectFactory(ClockColorEffectFactory)})
}

func ClockColorEffectFactory(l lampbase.WordClock) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		conf, ok := config.(*ClockColorConfig)
		if !ok {
			return errors.New("Not a ClockColorConfig")
		}
		m := color.RGBAModel
		return l.Color(m.Convert(conf.Color).(color.RGBA))
	})
}
