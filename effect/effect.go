package effect

import (
	"fmt"
	"lamp/lampbase"
	"time"
)

type Effect interface {
	Apply() (time.Duration, error)
}

type Configurer interface {
	Configure(conf Config)
}

type Config interface{}

type Info struct {
	Name        string
	Description string
}
type ExtendedInfo struct {
	Info
	ConfigFactory func() Config
	Factory       interface{}
}

type DeviceEffectFactory func(p lampbase.Device) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect

func (e *ExtendedInfo) Compatible(lamp lampbase.Device) bool {
	switch fac := e.Factory.(type) {
	case DeviceEffectFactory:
		_, ok := lamp.(lampbase.Device)
		return ok
	case DimLampEffectFactory:
		_, ok := lamp.(lampbase.DimLamp)
		return ok
	case ColorLampEffectFactory:
		_, ok := lamp.(lampbase.ColorLamp)
		return ok
	case StripeLampEffectFactory:
		_, ok := lamp.(lampbase.StripeLamp)
		return ok
	default:
		panic("Unknow lamp factory type " + fmt.Sprint(fac))
	}
	return false
}
