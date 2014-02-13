package effect

import (
	"log"
	"time"
)

type Controller struct {
	EffectChan chan Effect
	effect     Effect
}

func NewController() *Controller {
	return &Controller{EffectChan: make(chan Effect, 1)}
}

func (f *Controller) applyEffect() <-chan time.Time {
	var wait <-chan time.Time
	if dur, err := f.effect.Apply(); err != nil {
		log.Println(err)
		// Kill old effect and wait for new one
		f.effect = nil
		wait = nil
	} else if dur >= 0 {
		wait = time.After(dur)
	} else {
		wait = nil
	}
	return wait
}

func (f *Controller) Run() {
	var wait <-chan time.Time
	for {
		select {
		case f.effect = <-f.EffectChan:
			wait = f.applyEffect()
		case <-wait:
			wait = f.applyEffect()
		}
	}

}
