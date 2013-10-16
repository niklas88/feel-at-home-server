package wheel

import (
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type WheelConfig struct {
	Delay string
}

type WheelAll struct {
	wheelPos uint32
	forward  bool
	lamp     lampbase.ColorLamp
	delay    time.Duration
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
			Name:        "Wheel",
			Description: "A color wheel effect for color lamps"},
		ConfigFactory: func() effect.Config { return &WheelConfig{"30ms"} },
		Factory:       effect.ColorLampEffectFactory(NewWheelAllEffect)})
}

func NewWheelAllEffect(l lampbase.ColorLamp) effect.Effect {
	return &WheelAll{0, false, l, 30 * time.Millisecond}
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
	var c color.RGBA
	c.R, c.G, c.B = wheelColor(uint8(w.wheelPos))
	if w.wheelPos <= 0 || w.wheelPos >= 255 {
		w.forward = !w.forward
	}

	if w.forward {
		w.wheelPos++
	} else {
		w.wheelPos--
	}
	err := w.lamp.SetColor(&c)
	return w.delay, err
}
