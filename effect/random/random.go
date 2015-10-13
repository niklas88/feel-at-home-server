package random

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Random Brightness",
			Description: "Sets random pixels to random brightness"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.StripeLampEffectFactory(NewRandomBrightnessEffect)})

	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Random Fade",
			Description: "Fades randomly selected pixels"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.StripeLampEffectFactory(NewRandomFadeEffect)})

	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Random Color",
			Description: "Sets random pixels to random colors"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.StripeLampEffectFactory(NewRandomColorEffect)})

}

func NewRandomBrightnessEffect(l lampbase.StripeLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		sunriseConf, ok := config.(*effect.DelayConfig)
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

func NewRandomFadeEffect(l lampbase.StripeLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		sunriseConf, ok := config.(*effect.DelayConfig)
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

func NewRandomColorEffect(l lampbase.StripeLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		sunriseConf, ok := config.(*effect.DelayConfig)
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
