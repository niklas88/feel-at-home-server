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
	Effect        Effect
	ConfigFactory func() Config
}

func (e *Registration) Compatible(lamp lampbase.Device) bool {
	switch eff := e.Effect.(type) {
	case DeviceEffect:
		_, ok := lamp.(lampbase.Device)
		return ok
	case DimLampEffect:
		_, ok := lamp.(lampbase.DimLamp)
		return ok
	case ColorLampEffect:
		_, ok := lamp.(lampbase.ColorLamp)
		return ok
	case StripeLampEffect:
		_, ok := lamp.(lampbase.StripeLamp)
		return ok
	default:
		panic("Unknow effect type " + fmt.Sprintf("%g", eff))
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

func (r *Registry) ApplyEffect(name string, lamp lampbase.Device, config Config) (error, *Info) {
	r.RLock()
	defer r.RUnlock()
	e, ok := r.r[name]
	if !ok {
		return fmt.Errorf("Unknown Effect %s", name), nil
	}

	var err error

	switch eff := e.Effect.(type) {
	case DeviceEffect:
		err = eff(lamp, config)
	case DimLampEffect:
		err = eff(lamp.(lampbase.DimLamp), config)
	case ColorLampEffect:
		err = eff(lamp.(lampbase.ColorLamp), config)
	case StripeLampEffect:
		err = eff(lamp.(lampbase.StripeLamp), config)
	default:
		panic("Unknow effect type " + fmt.Sprint(eff))
	}
	return err, &e.Info
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
