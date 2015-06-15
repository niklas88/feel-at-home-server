package lampbase

import (
	"errors"
	"image/color"
	"time"
)

type UdpColorLamp struct {
	UdpDimLamp
}

func NewUdpColorLamp() *UdpColorLamp {
	return new(UdpColorLamp)
}

func (l *UdpColorLamp) writeColor(col color.Color) {
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf.WriteByte(byte(c.R))
	l.buf.WriteByte(byte(c.G))
	l.buf.WriteByte(byte(c.B))
}

func (l *UdpColorLamp) Color(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('C', 0x00)
	l.writeColor(col)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) ColorFade(delay time.Duration, col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('C', 0x01)
	l.UdpDimLamp.writeDurationMilliseconds(delay)
	l.writeColor(col)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) Sunrise(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('C', 0x02)
	l.UdpDimLamp.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) ColorWheel(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('C', 0x03)
	l.UdpDimLamp.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}
