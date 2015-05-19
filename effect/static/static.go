package static

import (
	"errors"
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
)

type StaticConfig struct {
	Color hexcolor.Hex
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Static",
			Description: "Set a static color for your color lamp"},
		ConfigFactory: func() effect.Config {
			return &StaticConfig{"#ffffff"}
		},
		Effect: effect.ColorLampEffect(StaticEffect)})
}

func StaticEffect(l lampbase.ColorLamp, c effect.Config) error {
	config, ok := c.(*StaticConfig)
	if !ok {
		return errors.New("Not a StaticConfig")
	}
	m := color.RGBAModel
	return l.SetColor(m.Convert(config.Color).(color.RGBA))
}
