package sunrise

import (
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type SunriseConfig struct {
	Delay string
}

type Sunrise struct {
	lamp    lampbase.ColorLamp
	delay   time.Duration
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Sunrise",
			Description: "A sunrise effect for color lamps"},
		ConfigFactory: func() effect.Config { return &SunriseConfig{"30ms"} },
		Factory:       effect.ColorLampEffectFactory(NewSunriseEffect)})
}

func NewSunriseEffect(l lampbase.ColorLamp) effect.Effect {
	return &Sunrise{l, 30 * time.Millisecond}
}

func (e *Sunrise) Configure(conf effect.Config) {
	sunriseConf := conf.(*SunriseConfig)
	var err error

	e.delay, err = time.ParseDuration(sunriseConf.Delay)
	if err != nil {
		log.Println(err)
		e.delay = 30 * time.Millisecond
	}
}

func (e *Sunrise) Apply() (time.Duration, error) {
	err := e.lamp.Sunrise(e.delay)
	return -1, err
}
