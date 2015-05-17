package lampbase

import (
	"errors"
	"time"
)

type UdpStripeLamp struct {
	UdpColorLamp
}

func NewUdpStripeLamp(numStripes, ledsPerStripe int) *UdpStripeLamp {
	lamp := new(UdpStripeLamp)
	return lamp
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
