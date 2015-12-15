package clockcolor

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"github.com/pwaller/go-hexcolor"
	"image/color"
)

type ClockColorConfig struct {
	Color hexcolor.Hex
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Clock Color",
			Description: "Set the color for your clock after clock effect selected"},
		ConfigFactory: func() deviceapi.Config {
			return &ClockColorConfig{"#ffffff"}
		},
		EffectFactory: deviceapi.WordClockEffectFactory(ClockColorEffectFactory)})
}

func ClockColorEffectFactory(l devices.WordClock) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		conf, ok := config.(*ClockColorConfig)
		if !ok {
			return errors.New("Not a ClockColorConfig")
		}
		m := color.RGBAModel
		return l.ClockColor(m.Convert(conf.Color).(color.RGBA))
	})
}
