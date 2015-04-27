package lampbase

import (
	"bytes"
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

func writeColor(buf *bytes.Buffer, col color.Color) {
	c := color.RGBAModel.Convert(col).(color.RGBA)
	buf.WriteByte(byte(c.R))
	buf.WriteByte(byte(c.G))
	buf.WriteByte(byte(c.B))
}

func (l *UdpColorLamp) SetColor(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('C', 0x00)
	writeColor(&l.buf, col)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) Colorfade(delay time.Duration, col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('C', 0x01)
	l.UdpDimLamp.writeDurationMilliseconds(delay)
	writeColor(&l.buf, col)
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
