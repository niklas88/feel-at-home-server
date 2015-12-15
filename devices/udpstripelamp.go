package devices

import (
	"bytes"
	"errors"
	"time"
)

type UdpStripeLamp struct {
	UdpColorLamp
}

func NewUdpStripeLamp() *UdpStripeLamp {
	lamp := new(UdpStripeLamp)
	return lamp
}

func (l *UdpStripeLamp) RandomPixelBrightness(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('S', 0x00, &buf)
	l.UdpColorLamp.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) RandomPixelWhiteFade(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('S', 0x01, &buf)
	l.UdpColorLamp.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) RandomPixelColor(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('S', 0x02, &buf)
	l.UdpColorLamp.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpStripeLamp) Rainbow(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('S', 0x04, &buf)
	l.UdpColorLamp.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}
