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
	lamp    lampbase.DimLamp
	delay   time.Duration
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Whitefade",
			Description: "Fades with white color"},
		ConfigFactory: func() effect.Config { return &WhitefadeConfig{"15ms"} },
		Factory:       effect.DimLampEffectFactory(NewWhitefadeEffect)})
}

func NewWhitefadeEffect(l lampbase.DimLamp) effect.Effect {
	return &Whitefade{l, 15 * time.Millisecond}
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

	err := w.lamp.Fade(w.delay, 255)
	return -1, err

}
