package static

import (
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type StaticConfig struct {
	Color hexcolor.Hex
}

type StaticEffect struct {
	color color.RGBA
	lamp  lampbase.ColorLamp
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Static",
			Description: "Set a static color for your color lamp"},
		ConfigFactory: func() effect.Config {
			return &StaticConfig{"#ffffff"}
		},
		Factory: effect.ColorLampEffectFactory(NewStaticEffect)})
}

func NewStaticEffect(l lampbase.ColorLamp) effect.Effect {
	return &StaticEffect{lamp: l}
}

func (s *StaticEffect) Configure(c effect.Config) {
	config := c.(*StaticConfig)
	m := color.RGBAModel
	s.color = m.Convert(config.Color).(color.RGBA)
}

func (s *StaticEffect) Apply() (time.Duration, error) {
	err := s.lamp.SetColor(s.color)
	return -1, err
}
