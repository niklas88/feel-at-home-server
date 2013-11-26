package sunrise

import (
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type SunriseConfig struct {
	Delay string
}

type Sunrise struct {
	step    uint32
	current color.RGBA
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
	return &Sunrise{0, color.RGBA{0, 0, 0, 0}, l, 30 * time.Millisecond}
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
	var r, g, b float64
	if e.step < 5000 {
		r = 204.0 * float64(e.step) / 5000.0
	} else if e.step < 10000 {
		r = 204.0 + (255.0-205.0)*float64(e.step-5000)/5000.0
		g = 51.0 * float64(e.step-5000) / 5000.0
	} else if e.step < 15000 {
		r = 255.0
		g = 51.0 + (255.0-51.0)*float64(e.step-10000)/5000.0
		b = 0
	} else if e.step < 20000 {
		r = 255.0
		g = 255.0
		b = 255.0 * float64(e.step-15000) / 5000.0
	} else if e.step >= 20000 {
		r = 255.0
		g = 255.0
		b = 255.0
		e.delay = -1
	}
	e.step++

	e.current.R, e.current.G, e.current.B = uint8(r), uint8(g), uint8(b)
	err := e.lamp.SetColor(&e.current)
	return e.delay, err
}
