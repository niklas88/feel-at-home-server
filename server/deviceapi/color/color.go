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
	deviceapi.DefaultRegistry.Register(deviceapi.NewColorLampEffect(
		"Color",
		"Set a static color four your lamp",
		applyToDevice,
		func() deviceapi.Config { return &ColorConfig{"#ffffff"} }))
}

func applyToDevice(l devices.ColorLamp, config deviceapi.Config) error {
	conf, ok := config.(*ColorConfig)
	if !ok {
		return errors.New("Not a ColorConfig")
	}
	m := color.RGBAModel
	return l.Color(m.Convert(conf.Color).(color.RGBA))
}
