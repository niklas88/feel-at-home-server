package whitefade

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"time"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Whitefade",
			Description: "White fading deviceapi"},
		ConfigFactory: deviceapi.DelayConfigFactory,
		EffectFactory: deviceapi.DimLampEffectFactory(NewWhiteFadeEffect)})
}

func NewWhiteFadeEffect(l devices.DimLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		strobeConf, ok := config.(*deviceapi.DelayConfig)
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
