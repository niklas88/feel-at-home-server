package random

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Random Brightness",
			Description: "Sets random pixels to random brightness"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.StripeLampEffectFactory(NewRandomBrightnessEffect)})

	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Random Fade",
			Description: "Fades randomly selected pixels"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.StripeLampEffectFactory(NewRandomFadeEffect)})

	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Random Color",
			Description: "Sets random pixels to random colors"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.StripeLampEffectFactory(NewRandomColorEffect)})

}

func NewRandomBrightnessEffect(l devices.StripeLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		sunriseConf, ok := config.(*deviceapi.DelayConfig)
		if !ok {
			return errors.New("Not a WheelConfig")
		}

		delay, err := time.ParseDuration(sunriseConf.Delay)
		if err != nil {
			return err
		}
		return l.RandomPixelBrightness(delay)
	})
}

func NewRandomFadeEffect(l devices.StripeLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		sunriseConf, ok := config.(*deviceapi.DelayConfig)
		if !ok {
			return errors.New("Not a WheelConfig")
		}

		delay, err := time.ParseDuration(sunriseConf.Delay)
		if err != nil {
			return err
		}
		return l.RandomPixelWhiteFade(delay)
	})
}

func NewRandomColorEffect(l devices.StripeLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		sunriseConf, ok := config.(*deviceapi.DelayConfig)
		if !ok {
			return errors.New("Not a WheelConfig")
		}

		delay, err := time.ParseDuration(sunriseConf.Delay)
		if err != nil {
			return err
		}
		return l.RandomPixelColor(delay)
	})
}
