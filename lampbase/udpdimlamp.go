package lampbase

import (
	"bytes"
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
	var buf bytes.Buffer
	l.UdpPowerDevice.writeHead('D', 0x00, &buf)
	buf.WriteByte(byte(b))
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) BrightnessScaling(b uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpPowerDevice.writeHead('D', 0x03, &buf)
	buf.WriteByte(byte(b))
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) writeDurationMilliseconds(delay time.Duration, buf *bytes.Buffer) {
	delayMilli := delay / time.Millisecond
	if delayMilli > math.MaxUint32 {
		delayMilli = math.MaxUint32
	}
	delaySmall := uint32(delayMilli)
	buf.WriteByte(byte(delaySmall >> 24))
	buf.WriteByte(byte((delaySmall >> 16) & 0xff))
	buf.WriteByte(byte((delaySmall >> 8) & 0xff))
	buf.WriteByte(byte(delaySmall & 0xff))
}

func (l *UdpDimLamp) Fade(delay time.Duration, maxBrightness uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpPowerDevice.writeHead('D', 0x01, &buf)
	l.writeDurationMilliseconds(delay, &buf)
	buf.WriteByte(byte(maxBrightness))
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) Stroboscope(delay time.Duration) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpPowerDevice.writeHead('D', 0x02, &buf)
	l.writeDurationMilliseconds(delay, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}
