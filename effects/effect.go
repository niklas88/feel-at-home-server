package effects

import (
	"lamp/lampbase"
	"time"
)

type Effect interface {
	Apply() (time.Duration, error)
}

type ConfigurerEffect interface {
	Effect
	Configure(conf interface{})
}

type EffectInfo struct {
	Name          string
	ConfigFactory func() interface{}
	Factory       interface{}
}

type DeviceEffectFactory func(p lampbase.Device) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect

func (e *EffectInfo) CreateEffect(lamp lampbase.Device) Effect {

	switch fac := e.Factory.(type) {
	case DeviceEffectFactory:
		if l, ok := lamp.(lampbase.Device); ok {
			return fac(l)
		}
	case DimLampEffectFactory:
		if l, ok := lamp.(lampbase.DimLamp); ok {
			return fac(l)
		}
	case ColorLampEffectFactory:
		if l, ok := lamp.(lampbase.ColorLamp); ok {
			return fac(l)
		}
	case StripeLampEffectFactory:
		if l, ok := lamp.(lampbase.StripeLamp); ok {
			return fac(l)
		}
	default:
		panic("Unknow lamp factory type")
	}
	return nil
}
