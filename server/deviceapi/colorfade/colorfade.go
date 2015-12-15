package colorefade

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"time"
)

type ColorfadeConfig struct {
	Color hexcolor.Hex
	Delay string
}

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Colorfade",
			Description: "Fades with Color"},
		ConfigFactory: func() deviceapi.Config { return &ColorfadeConfig{"#ffffff", "15ms"} },
		EffectFactory: deviceapi.ColorLampEffectFactory(ColorFadeEffect)})
}

func ColorFadeEffect(l devices.ColorLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
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
