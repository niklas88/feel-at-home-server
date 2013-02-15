package effects

import (
	"lamp/lampbase"
	"launchpad.net/tomb"
)

type Effect interface {
	Config() interface{}
	ConfigChan() chan interface{}
	Tomb() *tomb.Tomb
	Apply()
}

type PowerableEffectFactory func(p lampbase.Powerable) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect
