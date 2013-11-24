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

func wheelColor(w uint8) (uint8, uint8, uint8) {
	if w < 85 {
		return w * 3, 255 - w*3, 0
	} else if w < 170 {
		w -= 85
		return 255 - w*3, 0, w * 3
	}
	w -= 170
	return 0, w * 3, 255 - w*3

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
	var c color.RGBA
	if e.step < 5000 {
		c.R = uint8(204.0 * float64(e.step/5000.0))
		c.G = 0
		c.B = 0
	} else if e.step < 10000 {
		c.R = uint8(255.0 * float64((e.step-5000)/5000.0))
		c.G = uint8(51.0 * float64((e.step-5000)/5000.0))
		c.B = 0
	} else if e.step < 15000 {
		c.R = uint8(255.0 * float64((e.step-10000)/5000.0))
		c.G = uint8(255.0 * float64((e.step-10000)/5000.0))
		c.B = 0
	} else if e.step < 20000 {
		c.R = uint8(255.0 * float64((e.step-15000)/5000.0))
		c.G = uint8(255.0 * float64((e.step-15000)/5000.0))
		c.B = uint8(255.0 * float64((e.step-15000)/5000.0))
	} else if e.step >= 20000 {
		c.R = 255
		c.G = 255
		c.B = 255
		e.delay = -1
	}
	e.step++

	err := e.lamp.SetColor(&c)
	return e.delay, err
}
