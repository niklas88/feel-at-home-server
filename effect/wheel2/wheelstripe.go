package wheel2

import (
	"log"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type WheelConfig struct {
	Delay string
}

type WheelStripe struct {
	lamp     lampbase.StripeLamp
	delay    time.Duration
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Wheel2",
			Description: "Color wheel effect that sets single leds"},
		ConfigFactory: func() effect.Config { return &WheelConfig{"30ms"} },
		Factory:       effect.StripeLampEffectFactory(NewWheel2Effect)})
}

func (f *WheelStripe) Apply() (time.Duration, error) {
	err := f.lamp.Rainbow(f.delay)
	return -1, err
}

func NewWheel2Effect(l lampbase.StripeLamp) effect.Effect {
	return &WheelStripe{l, 30 * time.Millisecond}
}

func (f *WheelStripe) Configure(conf effect.Config) {
	wheelConf := conf.(*WheelConfig)
	var err error

	f.delay, err = time.ParseDuration(wheelConf.Delay)
	if err != nil {
		log.Println(err)
		f.delay = 30 * time.Millisecond
	}
}
