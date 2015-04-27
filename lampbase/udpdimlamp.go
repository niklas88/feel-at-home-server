package lampbase

import (
	"errors"
)

type UdpDimLamp struct {
	UdpPowerDevice
}

func NewUdpDimLamp() *UdpDimLamp {
	return new(UdpDimLamp)
}

func (l *UdpDimLamp) SetBrightness(b uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf.Reset()
	l.buf.WriteByte(byte(l.devicePort))
	l.buf.WriteByte('D')
	l.buf.WriteByte(0x00)
	l.buf.WriteByte(byte(b))
	_, err := l.buf.WriteTo(l.trans)
	return err
}
