package effect

import (
	"lamp/lampbase"
)

type Config interface{}

type DelayConfig struct {
	Delay string
}

func DelayConfigFactory() Config {
	return &DelayConfig{"30ms"}
}

type Info struct {
	Name        string
	Description string
	Config      Config
}

type Effect interface {
	Apply(config Config) error
}

type EffectFunc func(config Config) error

func (f EffectFunc) Apply(config Config) error {
	return f(config)
}

type DeviceEffectFactory func(p lampbase.Device) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect
