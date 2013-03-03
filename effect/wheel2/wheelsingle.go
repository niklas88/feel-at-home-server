package wheel2

import (
	"lamp/effect"
	"lamp/lampbase"
	"time"
)

type WheelSingle struct {
	wheelPos uint32
	lamp     lampbase.StripeLamp
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Wheel2",
			Description: "Color wheel effect that sets single leds",
			Config:      nil},
		Factory: effect.StripeLampEffectFactory(NewWheel2Effect)})
}

func (f *WheelSingle) Apply() (time.Duration, error) {
	f.colorizeLamp()
	f.lamp.UpdateAll()
	return 60 * time.Millisecond, nil
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

func NewWheel2Effect(l lampbase.StripeLamp) effect.Effect {
	return &WheelSingle{0, l}
}

func (f *WheelSingle) colorizeLamp() {
	stripes := f.lamp.Stripes()
	for _, s := range stripes {
		for i := 0; i < len(s); i++ {
			r, g, b := wheelColor(uint8((f.wheelPos + uint32(i*5)) % 256))
			s[i].R, s[i].G, s[i].B = r, g, b
		}
	}
	f.wheelPos++
	f.lamp.UpdateAll()
}
