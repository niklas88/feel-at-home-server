package effect

import (
	"errors"
	"fmt"
	"lamp/lampbase"
	"sync"
)

type Registry struct {
	sync.RWMutex
	r map[string]*Registration
}

type Registration struct {
	Info          Info
	EffectFactory interface{}
	ConfigFactory func() Config
}

func (e *Registration) Compatible(lamp lampbase.Device) bool {
	switch fac := e.EffectFactory.(type) {
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
	case MatrixLampEffectFactory:
		_, ok := lamp.(lampbase.MatrixLamp)
		return ok
	case WordClockEffectFactory:
		_, ok := lamp.(lampbase.WordClock)
		return ok
	default:
		panic("Unknow effect type " + fmt.Sprintf("%g", fac))
	}
}

var DefaultRegistry Registry

func init() {
	DefaultRegistry.r = make(map[string]*Registration)
}

func (r *Registry) Register(reg *Registration) error {
	r.Lock()
	defer r.Unlock()
	info := reg.Info
	if _, ok := r.r[info.Name]; ok {
		return errors.New("Tried adding two effects under Name " + info.Name)
	}
	// Populate Config with default
	reg.Info.Config = reg.ConfigFactory()
	r.r[info.Name] = reg
	return nil
}

func (r *Registry) Effect(name string, device lampbase.Device) Effect {
	r.RLock()
	defer r.RUnlock()
	e, ok := r.r[name]
	if !ok {
		return nil
	}
	var eff Effect
	switch fac := e.EffectFactory.(type) {
	case DeviceEffectFactory:
		eff = fac(device)
	case DimLampEffectFactory:
		d, ok := device.(lampbase.DimLamp)
		if !ok {
			return nil
		}
		eff = fac(d)
	case ColorLampEffectFactory:
		d, ok := device.(lampbase.ColorLamp)
		if !ok {
			return nil
		}
		eff = fac(d)
	case StripeLampEffectFactory:
		d, ok := device.(lampbase.StripeLamp)
		if !ok {
			return nil
		}
		eff = fac(d)
	case MatrixLampEffectFactory:
		d, ok := device.(lampbase.MatrixLamp)
		if !ok {
			return nil
		}
		eff = fac(d)
	case WordClockEffectFactory:
		d, ok := device.(lampbase.WordClock)
		if !ok {
			return nil
		}
		eff = fac(d)
	default:
		panic("Unknow effect factory type " + fmt.Sprintf("%q", eff))
	}

	return eff
}

func (r *Registry) CompatibleEffects(lamp lampbase.Device) []Info {
	compatibles := make([]Info, 0, 10)
	r.RLock()
	defer r.RUnlock()
	for _, v := range r.r {
		if v.Compatible(lamp) {
			compatibles = append(compatibles, v.Info)
		}
	}
	return compatibles
}

func (r *Registry) Info(name string) *Info {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.r[name]
	if !ok {
		return nil
	}
	info := v.Info
	info.Config = v.ConfigFactory()
	return &info
}

func (r *Registry) Config(name string) Config {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.r[name]
	if !ok {
		return nil
	}
	return v.ConfigFactory()
}
