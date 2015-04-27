package lampbase

import (
	"errors"
	"net"
)

type UdpPowerDevice struct {
	trans      *ReliableUDPTransport
	devicePort uint8
	buf        []uint8
}

func NewUdpPowerDevice() *UdpPowerDevice {
	return &UdpPowerDevice{nil, 0, make([]uint8, 3)}
}

func (l *UdpPowerDevice) Power(on bool) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = l.devicePort
	l.buf[1] = 'P'
	if on {
		l.buf[2] = 1

	} else {
		l.buf[2] = 0
	}
	_, err := l.trans.Write(l.buf[:3])
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
