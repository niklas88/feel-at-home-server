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
	InfoFactory func() Info
	Factory     interface{}
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
	info := reg.InfoFactory()
	r.Lock()
	defer r.Unlock()
	if _, ok := r.r[info.Name]; ok {
		return errors.New("Tried adding two effects under Name " + info.Name)
	}
	r.r[info.Name] = reg
	return nil
}

func (r *Registry) CreateEffect(info *Info, lamp lampbase.Device) Effect {
	e, ok := r.r[info.Name]
	if !ok {
		return nil
	}
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

func (r *Registry) CompatibleEffects(lamp lampbase.Device) []Info {
	compatibles := make([]Info, 0, 10)
	r.RLock()
	defer r.RUnlock()
	for _, v := range r.r {
		if v.Compatible(lamp) {
			compatibles = append(compatibles, v.InfoFactory())
		}
	}
	return compatibles
}

func (r *Registry) EffectInfo(name string) (Info, bool) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.r[name]
	if !ok {
		return Info{}, false
	}
	return v.InfoFactory(), true
}
