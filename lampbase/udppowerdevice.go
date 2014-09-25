package lampbase

import (
	"errors"
	"net"
)

type UdpPowerDevice struct {
	raddr      *net.UDPAddr
	laddr      *net.UDPAddr
	conn       *net.UDPConn
	devicePort uint8
	buf        []uint8
}

func NewUdpPowerDevice() *UdpPowerDevice {
	return &UdpPowerDevice{nil, nil, nil, 0, make([]uint8, 3)}
}

func (l *UdpPowerDevice) Power(on bool) error {
	l.buf[0] = l.devicePort
	l.buf[1] = 'P'
	if on {
		l.buf[2] = 1

	} else {
		l.buf[2] = 0
	}
	written, err := l.conn.Write(l.buf[:3])
	if err == nil && written != 3 {
		err = errors.New("Couldn't write udp packet in one call")
	}
	return err
}

func (l *UdpPowerDevice) Close() error {
	err := l.conn.Close()
	l.conn = nil
	return err
}

func (l *UdpPowerDevice) Dial(laddr, raddr *net.UDPAddr, devicePort uint8) (err error) {
	l.raddr, l.laddr = raddr, laddr
	l.devicePort = devicePort
	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err == nil {
		l.conn = conn
	}
	return
}
