package effects

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
	Name          string
	ConfigFactory func() Config
	Factory       interface{}
}

type DeviceEffectFactory func(p lampbase.Device) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect

func (e *Info) CreateEffect(lamp lampbase.Device) Effect {

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

func (e *Info) Compatible(lamp lampbase.Device) bool {
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
		panic("Unknow lamp factory type "+fmt.Sprint(fac))
	}
	return false
}
