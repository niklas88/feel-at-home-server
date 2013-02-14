package effects

import (
	"lamp/lampbase"
)

type Effect interface {
	ColorizeLamp(lamp lampbase.StripeLamp)
}

type EffectFunc func(lamp lampbase.StripeLamp)

func (f EffectFunc) ColorizeLamp(lamp lampbase.StripeLamp) {
	f(lamp)
}
