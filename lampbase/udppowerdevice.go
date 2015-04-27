package lampbase

import (
	"bytes"
	"errors"
	"net"
)

type UdpPowerDevice struct {
	trans   *ReliableUDPTransport
	lampNum uint8
	buf     bytes.Buffer
}

func NewUdpPowerDevice() *UdpPowerDevice {
	return new(UdpPowerDevice)
}

func (l *UdpPowerDevice) writeHead(effectGroup, effectNum byte) {
	l.buf.Reset()
	l.buf.WriteByte(byte(l.lampNum))
	l.buf.WriteByte(effectGroup)
	l.buf.WriteByte(effectNum)
}

func (l *UdpPowerDevice) Power(on bool) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.writeHead('P', 0x00)

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

func (l *UdpPowerDevice) Dial(laddr, raddr *net.UDPAddr, lampNum uint8) (err error) {
	l.lampNum = lampNum
	trans, err := DialReliableUDPTransport(laddr, raddr)
	if err == nil {
		l.trans = trans
	}
	return
}
