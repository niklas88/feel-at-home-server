package lampbase

import (
	"bytes"
	"errors"
	"net"
)

type UdpPowerDevice struct {
	trans      *ReliableUDPTransport
	devicePort uint8
	buf        bytes.Buffer
}

func NewUdpPowerDevice() *UdpPowerDevice {
	return new(UdpPowerDevice)
}

func (l *UdpPowerDevice) Power(on bool) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf.Reset()
	l.buf.WriteByte(byte(l.devicePort))
	l.buf.WriteByte('P')
	l.buf.WriteByte(0x00)

	if on {
		l.buf.WriteByte(1)

	} else {
		l.buf.WriteByte(0)
	}
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpPowerDevice) Close() error {
	err := l.trans.Close()
	l.trans = nil
	return err
}

func (l *UdpPowerDevice) Dial(laddr, raddr *net.UDPAddr, devicePort uint8) (err error) {
	l.devicePort = devicePort
	trans, err := DialReliableUDPTransport(laddr, raddr)
	if err == nil {
		l.trans = trans
	}
	return
}
