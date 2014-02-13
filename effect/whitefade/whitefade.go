package whitefade

import (
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"math"
	"time"
)

type WhitefadeConfig struct {
	Delay string
}

type Whitefade struct {
	current uint8
	upward  bool
	lamp    lampbase.DimLamp
	delay   time.Duration
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Whitefade",
			Description: "Fades with white color"},
		ConfigFactory: func() effect.Config { return &WhitefadeConfig{"30ms"} },
		Factory:       effect.DimLampEffectFactory(NewWhitefadeEffect)})
}

func NewWhitefadeEffect(l lampbase.DimLamp) effect.Effect {
	return &Whitefade{0, true, l, 30 * time.Millisecond}
}

func (w *Whitefade) Configure(conf effect.Config) {
	whitefadeConf := conf.(*WhitefadeConfig)
	var err error

	w.delay, err = time.ParseDuration(whitefadeConf.Delay)
	if err != nil {
		log.Println(err)
		w.delay = 30 * time.Millisecond
	}
}

func (w *Whitefade) Apply() (time.Duration, error) {

	err := w.lamp.SetBrightness(uint8(math.Pow(float64(w.current)/255, 2.5) * 255))
	if w.upward {
		w.current++
	} else {
		w.current--
	}

	if w.current == 255 || w.current == 0 {
		w.upward = !w.upward
	}
	return w.delay, err

}
