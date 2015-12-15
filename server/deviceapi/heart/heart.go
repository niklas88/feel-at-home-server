package heart

import (
	"errors"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
)

func init() {
	deviceapi.DefaultRegistry.Register(&deviceapi.Registration{
		Info: deviceapi.Info{
			Name:        "Heart",
			Description: "Set device into heart mode"},
		ConfigFactory: deviceapi.EmptyConfigFactory,
		EffectFactory: deviceapi.MatrixLampEffectFactory(HeartEffectFactory)})
}

func HeartEffectFactory(l devices.MatrixLamp) deviceapi.Effect {
	return deviceapi.EffectFunc(func(config deviceapi.Config) error {
		_, ok := config.(*deviceapi.EmptyConfig)
		if !ok {
			return errors.New("Not an empty Config")
		}

		return l.Heart()
	})
}
