package color

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"github.com/pwaller/go-hexcolor"
	"image/color"
)

type ColorConfig struct {
	Color hexcolor.Hex
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Color",
			Description: "Set a static color for your color lamp"},
		ConfigFactory: func() deviceapi.Config {
			return &ColorConfig{"#ffffff"}
		},
		EffectFactory: deviceapi.ColorLampEffectFactory(ColorEffectFactory)})
}

func ColorEffectFactory(l devices.ColorLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		conf, ok := config.(*ColorConfig)
		if !ok {
			return errors.New("Not a ColorConfig")
		}
		m := color.RGBAModel
		return l.Color(m.Convert(conf.Color).(color.RGBA))
	})
}
