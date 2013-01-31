package main

import (
	"fmt"
	"lamp/lampbase"
	"math/rand"
	"net"
	"time"
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

type borders struct {
	top float64
	bottom float64
}

func (bs *borders) Reset(r *rand.Rand, leds int) {
	desiredStdDev := float64(leds) * 0.07
	bs.top = r.NormFloat64()*desiredStdDev + float64(leds) * 0.80
	bs.bottom = r.NormFloat64()*desiredStdDev + float64(leds) * 0.40
}

func smooth(s lampbase.Stripe){
	o := make(lampbase.Stripe, len(s))
	copy(o, s)
	for i := 1; i < len(s)-2; i++ {
		s[i].R = uint8((float64(o[i-1].R)+2.0*float64(o[i].R)+float64(o[i+1].R))/4.0)
		s[i].G = uint8((float64(o[i-1].G)+2.0*float64(o[i].G)+float64(o[i+1].G))/4.0)
		s[i].B = uint8((float64(o[i-1].B)+2.0*float64(o[i].B)+float64(o[i+1].B))/4.0)
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "192.168.0.177:8888")
	if err != nil {
		fmt.Println("Couldn't resolve", err)
	}
	lamp := lampbase.NewLamp(4, 26, addr)

	r := rand.New(rand.NewSource(42))
	lamp.Update()
	borders := make([]borders, len(lamp.Stripes))

	desiredStdDev := float64(len(lamp.Stripes[0])) * 0.04

	for i, _ := range borders {
		borders[i].Reset(r, len(lamp.Stripes[i]))
	}

	for true {
		for strpn, s := range lamp.Stripes {
			borders[strpn].top += r.NormFloat64()*desiredStdDev
			borders[strpn].bottom += r.NormFloat64()*desiredStdDev
			bottom := clamp(borders[strpn].bottom, 0, len(s)-1)
			top := clamp(borders[strpn].top, 0, len(s)-1)

			for i := 0; i < bottom; i++ {
				s[i].R, s[i].G, s[i].B = 217, 93, 0
			}
			for i := bottom; i < top; i++ {
				s[i].R, s[i].G, s[i].B = 255, 0, 0
			}
			for i := top; i < len(s); i++ {
				s[i].R, s[i].G, s[i].B = 0, 0, 0
			}
			smooth(s)
		}
		lamp.Update()
		kill := r.Intn(300)
		if kill < len(borders) {
			borders[kill].Reset(r, len(lamp.Stripes[kill]))
		}
		time.Sleep(30 * time.Millisecond)
	}
}
