package effect

import (
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
	Config      Config
}

type DeviceEffectFactory func(p lampbase.Device) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect
