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
	Factory       interface{}
	ConfigFactory func() Config
}

func (e *Registration) Compatible(lamp lampbase.Device) bool {
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

func (r *Registry) CreateEffect(name string, lamp lampbase.Device) (Effect, *Info) {
	r.RLock()
	defer r.RUnlock()
	e, ok := r.r[name]
	if !ok {
		return nil, nil
	}
	var (
		eff  Effect
		info Info
	)
	switch fac := e.Factory.(type) {
	case DeviceEffectFactory:
		if l, ok := lamp.(lampbase.Device); ok {
			eff, info = fac(l), e.Info
		}
	case DimLampEffectFactory:
		if l, ok := lamp.(lampbase.DimLamp); ok {
			eff, info = fac(l), e.Info
		}
	case ColorLampEffectFactory:
		if l, ok := lamp.(lampbase.ColorLamp); ok {
			eff, info = fac(l), e.Info
		}
	case StripeLampEffectFactory:
		if l, ok := lamp.(lampbase.StripeLamp); ok {
			eff, info = fac(l), e.Info
		}
	default:
		panic("Unknow lamp factory type")
	}
	return eff, &info
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

func (r *Registry) Info(name string) (*Info, bool) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.r[name]
	if !ok {
		return nil, false
	}
	info := v.Info
	info.Config = v.ConfigFactory()
	return &info, true
}

func (r *Registry) Config(name string) (Config, bool) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.r[name]
	if !ok {
		return nil, false
	}
	return v.ConfigFactory(), true
}
