package wheel

import (
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type WheelConfig struct {
	Delay string
}

type WheelAll struct {
	lamp     lampbase.ColorLamp
	delay    time.Duration
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Wheel",
			Description: "A color wheel effect for color lamps"},
		ConfigFactory: func() effect.Config { return &WheelConfig{"30ms"} },
		Factory:       effect.ColorLampEffectFactory(NewWheelAllEffect)})
}

func NewWheelAllEffect(l lampbase.ColorLamp) effect.Effect {
	return &WheelAll{l, 30 * time.Millisecond}
}

func (f *WheelAll) Configure(conf effect.Config) {
	wheelConf := conf.(*WheelConfig)
	var err error

	f.delay, err = time.ParseDuration(wheelConf.Delay)
	if err != nil {
		log.Println(err)
		f.delay = 30 * time.Millisecond
	}
}

func (w *WheelAll) Apply() (time.Duration, error) {
	err := w.lamp.ColorWheel(w.delay)
	return -1, err
}
