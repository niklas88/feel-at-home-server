package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type UdpAnalogColorLamp struct {
	raddr *net.UDPAddr
	laddr *net.UDPAddr
	conn  *net.UDPConn
	buf   []uint8
}

func NewUdpAnalogColorLamp() *UdpAnalogColorLamp {
	return &UdpAnalogColorLamp{nil, nil, nil, make([]uint8, 4)}
}

func (l *UdpAnalogColorLamp) Power(on bool) error {
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

func (l *UdpAnalogColorLamp) SetBrightness(b uint8) error {
	color := color.RGBA{b, b, b, 0}
	return l.SetColor(&color)
}

func (l *UdpAnalogColorLamp) SetColor(col color.Color) error {
	if l.conn == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = 'C'
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf[1], l.buf[2], l.buf[3] = c.R, c.G, c.B
	written, err := l.conn.Write(l.buf[:4])
	if err == nil && written != 4 {
		err = errors.New("Couldn't write udp packet in one call")
	}
	// Change internal model
	if err == nil {
		for _, stripe := range l.stripes {
			for i := range stripe {
				stripe[i] = c
			}
		}
	}
	return err
}

func (l *UdpAnalogColorLamp) Close() error {
	err := l.conn.Close()
	l.conn = nil
	return err
}

func (l *UdpAnalogColorLamp) Dial(laddr, raddr *net.UDPAddr) (err error) {
	l.raddr, l.laddr = raddr, laddr

	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err == nil {
		l.conn = conn
	}
	return
}
