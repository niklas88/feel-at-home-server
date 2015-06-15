package lampbase

import (
	"errors"
	"math"
	"time"
)

type UdpDimLamp struct {
	UdpPowerDevice
}

func NewUdpDimLamp() *UdpDimLamp {
	return new(UdpDimLamp)
}

func (l *UdpDimLamp) Brightness(b uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpPowerDevice.writeHead('D', 0x00)
	l.buf.WriteByte(byte(b))
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) BrightnessScaling(b uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpPowerDevice.writeHead('D', 0x03)
	l.buf.WriteByte(byte(b))
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) writeDurationMilliseconds(delay time.Duration) {
	delayMilli := delay / time.Millisecond
	if delayMilli > math.MaxUint32 {
		delayMilli = math.MaxUint32
	}
	delaySmall := uint32(delayMilli)
	l.buf.WriteByte(byte(delaySmall >> 24))
	l.buf.WriteByte(byte((delaySmall >> 16) & 0xff))
	l.buf.WriteByte(byte((delaySmall >> 8) & 0xff))
	l.buf.WriteByte(byte(delaySmall & 0xff))
}

func (l *UdpDimLamp) Fade(delay time.Duration, maxBrightness uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpPowerDevice.writeHead('D', 0x01)
	l.writeDurationMilliseconds(delay)
	l.buf.WriteByte(byte(maxBrightness))
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) Stroboscope(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpPowerDevice.writeHead('D', 0x02)
	l.writeDurationMilliseconds(delay)
	_, err := l.buf.WriteTo(l.trans)
	return err
}
