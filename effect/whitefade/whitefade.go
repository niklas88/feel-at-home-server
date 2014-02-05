package whitefade

import (
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type WhitefadeConfig struct {
	Delay string
}

type Whitefade struct {
	step    uint8
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
	return &Whitefade{255, 255, true, l, 30 * time.Millisecond}
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

	if w.step < 255 {
		w.step++
	} else {
		w.step = 0
		w.upward = !w.upward
	}
	if w.upward {
		w.current++
	} else {
		w.current--
	}

	err := w.lamp.SetBrightness(w.current)
	return w.delay, err
}
