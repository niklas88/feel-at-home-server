package lampbase

import (
	"errors"
	"net"
)

type UdpDimLamp struct {
   trans      *ReliableUDPTransport
	devicePort uint8
	buf        []uint8
}

func NewUdpDimLamp() *UdpDimLamp {

	return &UdpDimLamp{nil, 0, make([]uint8, 3)}
}

func (l *UdpDimLamp) Power(on bool) error {
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
	err :=  l.trans.SendReliable(l.buf[:3])
	return err
}

func (l *UdpDimLamp) SetBrightness(b uint8) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = l.devicePort
	l.buf[1] = 'B'
	l.buf[2] = b
   err :=  l.trans.SendReliable(l.buf[:3])
	return err
}

func (l *UdpDimLamp) Close() error {
	err := l.trans.Close()
	l.trans = nil
	return err
}

func (l *UdpDimLamp) Dial(laddr, raddr *net.UDPAddr, devicePort uint8) (err error) {
	l.devicePort = devicePort

	trans, err := DialReliableUDPTransport(laddr, raddr)
	if err == nil {
		l.trans = trans
	}
	return
}
