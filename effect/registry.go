package effect

import (
	"errors"
	"lamp/lampbase"
	"sync"
)

type Registry struct {
	r map[string]*ExtendedInfo
	m sync.RWMutex
}

var DefaultRegistry Registry

func init() {
	DefaultRegistry.r = make(map[string]*ExtendedInfo)
}

func (r *Registry) Register(info *ExtendedInfo) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.r[info.Name]; ok {
		return errors.New("Tried adding two effects under Name " + info.Name)
	}
	r.r[info.Name] = info
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

func (r *Registry) CreateConfig(info *Info) Config {
	r.m.RLock()
	defer r.m.RUnlock()
	e, ok := r.r[info.Name]
	if !ok {
		return nil
	}
	return e.ConfigFactory()
}

func (r *Registry) CompatibleEffects(lamp lampbase.Device) []Info {
	compatibles := make([]Info, 0, 10)
	r.m.RLock()
	defer r.m.RUnlock()
	for _, v := range r.r {
		if v.Compatible(lamp) {
			compatibles = append(compatibles, Info{v.Name, v.Description})
		}
	}
	return compatibles
}
