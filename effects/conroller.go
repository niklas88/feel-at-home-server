package effects

import (
	"launchpad.net/tomb"
	"log"
	"time"
)

type Controller struct {
	EffectChan chan Effect
	Tomb       tomb.Tomb
	effect     Effect
}

func NewController() *Controller {
	return &Controller{EffectChan: make(chan Effect, 1)}
}

func (f *Controller) Run() {
	defer f.Tomb.Done()
	var wait <-chan time.Time
	for {
		select {
		case <-f.Tomb.Dying():
			close(f.EffectChan)
			return
		case f.effect = <-f.EffectChan:
			wait = time.After(0)
		case <-wait:
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
		}
	}

}
