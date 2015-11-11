package color

import (
	"errors"
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"github.com/niklas88/feel-at-home-server/effect"
	"github.com/niklas88/feel-at-home-server/lampbase"
)

type ColorConfig struct {
	Color hexcolor.Hex
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Color",
			Description: "Set a static color for your color lamp"},
		ConfigFactory: func() effect.Config {
			return &ColorConfig{"#ffffff"}
		},
		EffectFactory: effect.ColorLampEffectFactory(ColorEffectFactory)})
}

func ColorEffectFactory(l lampbase.ColorLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		conf, ok := config.(*ColorConfig)
		if !ok {
			return errors.New("Not a ColorConfig")
		}
		m := color.RGBAModel
		return l.Color(m.Convert(conf.Color).(color.RGBA))
	})
}
