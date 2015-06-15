package lampbase

import (
	"errors"
)

type UdpMatrixLamp struct {
	UdpStripeLamp
}

func NewUdpMatrixLamp() *UdpMatrixLamp {
	lamp := new(UdpMatrixLamp)
	return lamp
}

func (l *UdpMatrixLamp) Heart() error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('M', 0x00)
	_, err := l.buf.WriteTo(l.trans)
	return err
}
