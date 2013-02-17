package effects

import (
	"lamp/lampbase"
	"time"
)

type Effect interface {
	Apply() (time.Duration, error)
}

type EffectInfo struct {
	Name          string
	ConfigFactory func() interface{}
	Factory       interface{}
}

type PowerableEffectFactory func(p lampbase.Powerable, conf interface{}) Effect
type DimLampEffectFactory func(d lampbase.DimLamp, conf interface{}) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp, conf interface{}) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp, conf interface{}) Effect
