package effects

import (
	"image/color"
	"lamp/lampbase"
	"math/rand"
	"time"
)

type FireConfig struct {
	BottomColor color.RGBA
	MidColor    color.RGBA
	TopColor    color.RGBA
}

func clamp(val float64, lower, upper int) (ret int) {
	ret = int(val)
	if ret > upper {
		ret = upper
	} else if ret < lower {
		ret = lower
	}
	return
}

type borderpair struct {
	top    float64
	bottom float64
}

type FireEffect struct {
	r       *rand.Rand
	lamp    lampbase.StripeLamp
	borders []borderpair
	stdDev  float64
	config  *FireConfig
}

func (f *FireEffect) Apply() (time.Duration, error) {
	f.colorizeLamp()
	f.lamp.UpdateAll()
	return 30 * time.Millisecond, nil
}

func (f *FireEffect) colorizeLamp() {
	stripes := f.lamp.Stripes()
	for strpn, s := range stripes {
		f.borders[strpn].top += f.r.NormFloat64() * f.stdDev
		f.borders[strpn].bottom += f.r.NormFloat64() * f.stdDev
		bottom := clamp(f.borders[strpn].bottom, 0, len(s)-1)
		top := clamp(f.borders[strpn].top, 0, len(s)-1)

		for i := 0; i < bottom; i++ {
			s[i] = f.config.BottomColor //217, 93, 0
		}
		for i := bottom; i < top; i++ {
			s[i] = f.config.MidColor // 255, 0, 0
		}
		for i := top; i < len(s); i++ {
			s[i] = f.config.TopColor // 0,0,0
		}
		for i := 0; i < 5; i++ {
			smooth(s)
		}
	}
	kill := f.r.Intn(300)
	if kill < len(f.borders) {
		f.borders[kill].reset(f.r, len(stripes[kill]))
	}
	f.lamp.UpdateAll()
}

func (bs *borderpair) reset(r *rand.Rand, leds int) {
	desiredStdDev := float64(leds) * 0.04
	bs.top = r.NormFloat64()*desiredStdDev + float64(leds)*0.80
	bs.bottom = r.NormFloat64()*desiredStdDev + float64(leds)*0.30
}

func smooth(s lampbase.Stripe) {
	o := make(lampbase.Stripe, len(s))
	copy(o, s)
	for i := 1; i < len(s)-2; i++ {
		s[i].R = uint8((float64(o[i-1].R) + 2.0*float64(o[i].R) + float64(o[i+1].R)) / 4.0)
		s[i].G = uint8((float64(o[i-1].G) + 2.0*float64(o[i].G) + float64(o[i+1].G)) / 4.0)
		s[i].B = uint8((float64(o[i-1].B) + 2.0*float64(o[i].B) + float64(o[i+1].B)) / 4.0)
	}
}

func (f *FireEffect) Configure(conf interface{}) {
	f.config = conf.(*FireConfig)
}

func NewFireEffect(l lampbase.StripeLamp) Effect {
	stripes := l.Stripes()
	numStripes := len(stripes)
	f := &FireEffect{r: rand.New(rand.NewSource(42)), lamp: l, config: nil, borders: make([]borderpair, numStripes), stdDev: float64(len(stripes[0])) * 0.04}
	for i, b := range f.borders {
		b.reset(f.r, len(stripes[i]))
	}
	return f
}
