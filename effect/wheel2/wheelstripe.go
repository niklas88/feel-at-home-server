package wheel2

import (
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type WheelConfig struct {
	Delay string
}

type WheelStripe struct {
	wheelPos uint32
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
	f.colorizeLamp()
	err := f.lamp.UpdateAll()
	return f.delay, err
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
	return &WheelStripe{0, l, 30 * time.Millisecond}
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

func (f *WheelStripe) colorizeLamp() {
	stripes := f.lamp.Stripes()
	for _, s := range stripes {
		for i := 0; i < len(s); i++ {
			r, g, b := wheelColor(uint8((f.wheelPos + uint32(i*5)) % 256))
			s[i].R, s[i].G, s[i].B = r, g, b
		}
	}
	f.wheelPos--
}
