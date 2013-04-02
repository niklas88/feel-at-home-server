package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type UdpPowerDevice struct {
	raddr *net.UDPAddr
	laddr *net.UDPAddr
	conn  *net.UDPConn
	buf   []uint8
}

func NewPowerDevice() *UdpPowerDevice {
	return &UdpPowerDevice{nil, nil, nil, make([]uint8, 2)}
}

func (l *UdpPowerDevice) Power(on bool) error {
	l.buf[0] = 'P'
	if on {
		l.buf[1] = 1

	} else {
		l.buf[1] = 0
	}
	written, err := l.conn.Write(l.buf[:2])
	if err == nil && written != 2 {
		err = errors.New("Couldn't write udp packet in one call")
	}
	return err
}

func (l *UdpPowerDevice) Close() error {
	err := l.conn.Close()
	l.conn = nil
	return err
}

func (l *UdpPowerDevice) Dial(laddr, raddr *net.UDPAddr) (err error) {
	l.raddr, l.laddr = raddr, laddr

	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err == nil {
		l.conn = conn
	}
	return
}
