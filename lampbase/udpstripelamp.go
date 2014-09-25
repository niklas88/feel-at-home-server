package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type UdpStripeLamp struct {
	stripes    []Stripe
	raddr      *net.UDPAddr
	laddr      *net.UDPAddr
	conn       *net.UDPConn
	devicePort uint8
	buf        []uint8
}

func NewUdpStripeLamp(numStripes, ledsPerStripe int) *UdpStripeLamp {
	stripes := make([]Stripe, numStripes)
	for i := range stripes {
		stripes[i] = make(Stripe, ledsPerStripe)
	}
	return &UdpStripeLamp{stripes, nil, nil, nil, 0, make([]uint8, ledsPerStripe*numStripes*3+2)}
}

func (l *UdpStripeLamp) Power(on bool) error {
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

func (l *UdpStripeLamp) SetBrightness(b uint8) error {
	color := color.RGBA{b, b, b, 0}
	return l.SetColor(&color)
}

func (l *UdpStripeLamp) SetColor(col color.Color) error {
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

func (l *UdpStripeLamp) Stripes() []Stripe {
	return l.stripes
}

func (l *UdpStripeLamp) Close() error {
	err := l.conn.Close()
	l.conn = nil
	return err
}

func (l *UdpStripeLamp) Dial(laddr, raddr *net.UDPAddr, devicePort uint8) (err error) {
	l.raddr, l.laddr = raddr, laddr
	l.devicePort = devicePort

	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err == nil {
		l.conn = conn
	}
	return
}

func (l *UdpStripeLamp) UpdateAll() error {
	if l.conn == nil {
		return errors.New("Not Dialed")

	}
	l.buf[0] = l.devicePort
	l.buf[1] = 'D'
	bufpos := 2
	for i, s := range l.stripes {
		if i%2 == 0 {
			for j := 0; j < len(s); j++ {
				l.buf[bufpos], l.buf[bufpos+1], l.buf[bufpos+2] = s[j].R, s[j].G, s[j].B
				bufpos += 3
			}
		} else { // Go odd stripes in reverse
			for j := len(s) - 1; j >= 0; j-- {
				l.buf[bufpos], l.buf[bufpos+1], l.buf[bufpos+2] = s[j].R, s[j].G, s[j].B
				bufpos += 3
			}
		}
	}

	written, err := l.conn.Write(l.buf)
	if err == nil && written != len(l.buf) {
		err = errors.New("Couldn't write buf in single write/packet")
	}
	return err
}
