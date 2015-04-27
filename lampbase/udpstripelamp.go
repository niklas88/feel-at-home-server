package lampbase

import (
	"errors"
	"image/color"
	"time"
)

type UdpStripeLamp struct {
	UdpColorLamp
	stripes []Stripe
}

func NewUdpStripeLamp(numStripes, ledsPerStripe int) *UdpStripeLamp {
	stripes := make([]Stripe, numStripes)
	for i := range stripes {
		stripes[i] = make(Stripe, ledsPerStripe)
	}
	lamp := new(UdpStripeLamp)
	lamp.stripes = stripes
	return lamp
}

func (l *UdpStripeLamp) SetBrightness(b uint8) error {
	err := l.UdpColorLamp.SetBrightness(b)
	color := color.RGBA{b, b, b, 0}
	// Change model
	if err == nil {
		for _, stripe := range l.stripes {
			for i := range stripe {
				stripe[i] = color
			}
		}
	}
	return err
}

func (l *UdpStripeLamp) SetColor(col color.Color) error {
	err := l.UdpColorLamp.SetColor(col)
	// Change internal model
	if err == nil {
		for _, stripe := range l.stripes {
			for i := range stripe {
				stripe[i] = col.(color.RGBA)
			}
		}
	}
	return err
}

func (l *UdpStripeLamp) RandomPixelBrightness(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('S', 0x00)
	l.UdpColorLamp.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) RandomPixelWhiteFade(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('S', 0x01)
	l.UdpColorLamp.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) RandomPixelColor(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('S', 0x02)
	l.UdpColorLamp.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) Rainbow(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('S', 0x04)
	l.UdpColorLamp.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) Stripes() []Stripe {
	return l.stripes
}

func (l *UdpStripeLamp) UpdateAll() error {
	if l.trans == nil {
		return errors.New("Not Dialed")

	}
	l.UdpColorLamp.writeHead('S', 0x03)
	for i, s := range l.stripes {
		if i%2 == 0 {
			for j := 0; j < len(s); j++ {
				l.buf.WriteByte(byte(s[j].R))
				l.buf.WriteByte(byte(s[j].G))
				l.buf.WriteByte(byte(s[j].B))
			}
		} else { // Go odd stripes in reverse
			for j := len(s) - 1; j >= 0; j-- {
				l.buf.WriteByte(byte(s[j].R))
				l.buf.WriteByte(byte(s[j].G))
				l.buf.WriteByte(byte(s[j].B))
			}
		}
	}

	_, err := l.buf.WriteTo(l.trans)
	return err
}
