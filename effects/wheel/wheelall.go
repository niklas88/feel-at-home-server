package effect

import (
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type WheelAll struct {
	wheelPos uint32
	forward  bool
	lamp     lampbase.ColorLamp
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
	effect.DefaultRegistry.Register(&effect.ExtendedInfo{
		Info: effect.Info{
			Name:        "Wheel",
			Description: "A color wheel effect for color lamps"},
		ConfigFactory: func() effect.Config { return nil },
		Factory:       effect.ColorLampEffectFactory(NewWheelAllEffect)})
}

func NewWheelAllEffect(l lampbase.ColorLamp) effect.Effect {
	return &WheelAll{0, false, l}
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
	return 30 * time.Millisecond, err
}
