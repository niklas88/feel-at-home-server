package whitefade

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/effect"
	"github.com/niklas88/feel-at-home-server/lampbase"
	"time"
)

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Whitefade",
			Description: "White fading effect"},
		ConfigFactory: effect.DelayConfigFactory,
		EffectFactory: effect.DimLampEffectFactory(NewWhiteFadeEffect)})
}

func NewWhiteFadeEffect(l lampbase.DimLamp) effect.Effect {
	return effect.EffectFunc(func(config effect.Config) error {
		strobeConf, ok := config.(*effect.DelayConfig)
		if !ok {
			return errors.New("Not a WhiteFadeConfig")
		}
		delay, err := time.ParseDuration(strobeConf.Delay)
		if err != nil {
			return err
		}
		return l.Fade(delay, 255)
	})
}
