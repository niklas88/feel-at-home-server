package effects

import (
	"lamp/lampbase"
)

type Wheel struct {
	wheelPos uint32
	forward  bool
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

func (w *Wheel) ColorizeLamp(lamp *lampbase.Lamp) {
	for _, s := range lamp.Stripes {
		for i := range s {
			s[i].R, s[i].G, s[i].B = wheelColor(uint8(w.wheelPos))
		}
	}
	if w.wheelPos <= 0 || w.wheelPos >= 255 {
		w.forward = !w.forward
	}

	if w.forward {
		w.wheelPos++
	} else {
		w.wheelPos--
	}
}
