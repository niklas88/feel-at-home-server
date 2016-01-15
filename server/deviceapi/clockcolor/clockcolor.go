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

func applyToDevice(l devices.WordClock, config deviceapi.Config) error {
	conf, ok := config.(*ClockColorConfig)
	if !ok {
		return errors.New("Not a ClockColorConfig")
	}
	m := color.RGBAModel
	return l.ClockColor(m.Convert(conf.Color).(color.RGBA))
}

func init() {
	deviceapi.DefaultRegistry.Register(
		deviceapi.NewWordClockEffect(
			"Clock Color",
			"Set the color for your clock",
			applyToDevice,
			func() deviceapi.Config { return &ClockColorConfig{"#ffffff"} }))
}
