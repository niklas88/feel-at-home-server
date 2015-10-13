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

func (l *UdpColorLamp) writeColor(col color.Color, buf *bytes.Buffer) {
	c := color.RGBAModel.Convert(col).(color.RGBA)
	buf.WriteByte(byte(c.R))
	buf.WriteByte(byte(c.G))
	buf.WriteByte(byte(c.B))
}

func (l *UdpColorLamp) Color(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpDimLamp.writeHead('C', 0x00, &buf)
	l.writeColor(col, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) ColorFade(delay time.Duration, col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpDimLamp.writeHead('C', 0x01, &buf)
	l.UdpDimLamp.writeDurationMilliseconds(delay, &buf)
	l.writeColor(col, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) Sunrise(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpDimLamp.writeHead('C', 0x02, &buf)
	l.UdpDimLamp.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpColorLamp) ColorWheel(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpDimLamp.writeHead('C', 0x03, &buf)
	l.UdpDimLamp.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}
