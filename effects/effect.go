package effects

import (
	"lamp/lampbase"
)

type Effect interface {
	ColorizeLamp(lamp *lampbase.Lamp)
}

type EffectFunc func(lamp *lampbase.Lamp)

func (f EffectFunc) ColorizeLamp(lamp *lampbase.Lamp) {
	f(lamp)
}
