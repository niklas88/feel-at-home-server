package brightness

import (
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type BrightnessConfig struct {
	Brightness uint8
}

type BrightnessEffect struct {
	brightness uint8
	lamp       lampbase.DimLamp
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Brightness",
			Description: "Set brightness for your lamp"},
		ConfigFactory: func() effect.Config {
			return &BrightnessConfig{255}
		},
		Factory: effect.DimLampEffectFactory(NewBrightnessEffect)})
}

func NewBrightnessEffect(l lampbase.DimLamp) effect.Effect {
	return &BrightnessEffect{lamp: l}
}

func (b *BrightnessEffect) Configure(c effect.Config) {
	config := c.(*BrightnessConfig)
	b.brightness = config.Brightness
}

func (b *BrightnessEffect) Apply() (time.Duration, error) {
	err := b.lamp.SetBrightness(b.brightness)
	return -1, err
}
