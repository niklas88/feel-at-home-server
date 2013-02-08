package effects

import (
	"lamp/lampbase"
	"math/rand"
)

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
	lamp    *lampbase.Lamp
	borders []borderpair
	stdDev  float64
}

func (f *FireEffect) ColorizeLamp(lamp *lampbase.Lamp) {
	if f.lamp != lamp {
		f.setLamp(lamp)
	}

	for strpn, s := range lamp.Stripes {
		f.borders[strpn].top += f.r.NormFloat64() * f.stdDev
		f.borders[strpn].bottom += f.r.NormFloat64() * f.stdDev
		bottom := clamp(f.borders[strpn].bottom, 0, len(s)-1)
		top := clamp(f.borders[strpn].top, 0, len(s)-1)

		for i := 0; i < bottom; i++ {
			s[i].R, s[i].G, s[i].B = 217, 93, 0
		}
		for i := bottom; i < top; i++ {
			s[i].R, s[i].G, s[i].B = 255, 0, 0
		}
		for i := top; i < len(s); i++ {
			s[i].R, s[i].G, s[i].B = 0, 0, 0
		}
		for i := 0; i < 5; i++ {
			smooth(s)
		}
	}
	kill := f.r.Intn(300)
	if kill < len(f.borders) {
		f.borders[kill].Reset(f.r, len(lamp.Stripes[kill]))
	}
}

func (bs *borderpair) Reset(r *rand.Rand, leds int) {
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

func (f *FireEffect) setLamp(l *lampbase.Lamp) {
	f.lamp = l
	numStripes := len(l.Stripes)
	f.borders = make([]borderpair, numStripes)
	f.stdDev = float64(len(l.Stripes[0])) * 0.04

	for i, b := range f.borders {
		b.Reset(f.r, len(l.Stripes[i]))
	}
}

func NewFireEffect() *FireEffect {
	ret := &FireEffect{r: rand.New(rand.NewSource(42)), lamp: nil, borders: make([]borderpair, 0), stdDev: 0.0}
	return ret
}
