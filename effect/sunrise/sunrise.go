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
	current int
	lamp    lampbase.ColorLamp
	delay   time.Duration
}

type floatRGB struct {
	r float64
	g float64
	b float64
}

func (e *floatRGB) RGB() (r float64, g float64, b float64) {
	r = e.r
	g = e.g
	b = e.b
	return
}

const stepsPerColor = 2500

var colorList []floatRGB

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Sunrise",
			Description: "A sunrise effect for color lamps"},
		ConfigFactory: func() effect.Config { return &SunriseConfig{"30ms"} },
		Factory:       effect.ColorLampEffectFactory(NewSunriseEffect)})

	colorList = []floatRGB{
		floatRGB{0.0, 0.0, 0.0},
		floatRGB{0.0, 0.0, 10.0},
		floatRGB{10.0, 0.0, 5.0},
		floatRGB{50.0, 0.0, 0.0},
		floatRGB{100.0, 40.0, 0.0},
		floatRGB{160.0, 80.0, 0.0},
		floatRGB{200.0, 100.0, 0.0},
		floatRGB{255.0, 200.0, 100.0}}
}

func NewSunriseEffect(l lampbase.ColorLamp) effect.Effect {
	return &Sunrise{0, 0, l, 30 * time.Millisecond}
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
	if e.step >= stepsPerColor {
		e.step = 0
		e.current++
		if e.current >= len(colorList)-1 {
			e.delay = -1
			return e.delay, nil
		}
	}
	r, g, b := colorList[e.current].RGB()
	rn, gn, bn := colorList[e.current+1].RGB()
	mix := float64(e.step) / stepsPerColor
	r = r + mix*(rn-r)
	g = g + mix*(gn-g)
	b = b + mix*(bn-b)

	e.step++
	err := e.lamp.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), 0})
	return e.delay, err
}
