package lampbase

import (
	"errors"
	"image/color"
)

type UdpColorLamp struct {
	UdpDimLamp
}

func NewUdpColorLamp() *UdpColorLamp {
	return new(UdpColorLamp)
}

func (l *UdpColorLamp) SetColor(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf.Reset()
	l.buf.WriteByte(byte(l.devicePort))
	l.buf.WriteByte('C')
	l.buf.WriteByte(0x00)
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf.WriteByte(byte(c.R))
	l.buf.WriteByte(byte(c.G))
	l.buf.WriteByte(byte(c.B))
	_, err := l.buf.WriteTo(l.trans)
	return err
}
