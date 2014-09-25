package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type UdpAnalogColorLamp struct {
	raddr      *net.UDPAddr
	laddr      *net.UDPAddr
	conn       *net.UDPConn
	devicePort uint8
	buf        []uint8
}

func NewUdpAnalogColorLamp() *UdpAnalogColorLamp {
	return &UdpAnalogColorLamp{nil, nil, nil, 0, make([]uint8, 5)}
}

func (l *UdpAnalogColorLamp) Power(on bool) error {
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

func (l *UdpAnalogColorLamp) SetBrightness(b uint8) error {
	color := color.RGBA{b, b, b, 0}
	return l.SetColor(&color)
}

func (l *UdpAnalogColorLamp) SetColor(col color.Color) error {
	if l.conn == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = l.devicePort
	l.buf[1] = 'C'
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf[2], l.buf[3], l.buf[4] = c.R, c.G, c.B
	written, err := l.conn.Write(l.buf[:5])
	if err == nil && written != 5 {
		err = errors.New("Couldn't write udp packet in one call")
	}
	return err
}

func (l *UdpAnalogColorLamp) Close() error {
	err := l.conn.Close()
	l.conn = nil
	return err
}

func (l *UdpAnalogColorLamp) Dial(laddr, raddr *net.UDPAddr, devicePort uint8) (err error) {
	l.raddr, l.laddr = raddr, laddr
	l.devicePort = devicePort

	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err == nil {
		l.conn = conn
	}
	return
}
