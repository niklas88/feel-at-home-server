package random

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	configFactory := func() deviceapi.Config { return &deviceapi.DelayConfig{"30ms"} }
	deviceapi.DefaultRegistry.Register(deviceapi.NewStripeLampEffect(
		"Random Brightness",
		"Sets random pixels to randrom brightness",
		applyRandomBrightness,
		configFactory))

	deviceapi.DefaultRegistry.Register(deviceapi.NewStripeLampEffect(
		"Random Fade",
		"Fades randomly selected pixels",
		applyRandomFade,
		configFactory))

	deviceapi.DefaultRegistry.Register(deviceapi.NewStripeLampEffect(
		"Random Color",
		"Sets random pixels to random colors",
		applyRandomColor,
		configFactory))
}

func applyRandomBrightness(l devices.StripeLamp, config deviceapi.Config) error {
	sunriseConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a DelayConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.RandomPixelBrightness(delay)
}

func applyRandomFade(l devices.StripeLamp, config deviceapi.Config) error {
	sunriseConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a DelayConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.RandomPixelWhiteFade(delay)
}

func applyRandomColor(l devices.StripeLamp, config deviceapi.Config) error {
	sunriseConf, ok := config.(*deviceapi.DelayConfig)
	if !ok {
		return errors.New("Not a DelayConfig")
	}

	delay, err := time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		return err
	}
	return l.RandomPixelColor(delay)
}
