package devices

import (
	"bytes"
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
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('M', 0x00, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}
