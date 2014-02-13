package effect

import (
	"lamp/lampbase"
	"log"
	"time"
)

const (
	Activate   = iota
	Deactivate = iota
	Shutdown   = iota
)

type Controller struct {
	StateChange chan int
	EffectChan  chan Effect
	effect      Effect
	device      lampbase.Device
}

func NewController(device lampbase.Device) *Controller {
	return &Controller{EffectChan: make(chan Effect, 1), StateChange: make(chan int, 1), device: device}
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
	var (
		wait <-chan time.Time
	)
	for {
		select {
		case req := <-f.StateChange:
			switch req {
			case Activate:
				log.Println("Activate")
				if f.effect != nil {
					wait = f.applyEffect()
				} else {
					f.device.Power(true)
				}

			case Deactivate:
				log.Println("Deactivate")
				f.device.Power(false)
				wait = nil
			case Shutdown:
				log.Println("Shutdown")
				f.device.Power(false)
				wait = nil
				close(f.EffectChan)
				close(f.StateChange)
				return
			}
		case f.effect = <-f.EffectChan:
			wait = f.applyEffect()
		case <-wait:
			wait = f.applyEffect()
		}
	}

}
