package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type UdpStripeLamp struct {
	stripes []Stripe
	trans   *ReliableUDPTransport
	buf     []byte
	seqNum  uint8
}

func NewUdpStripeLamp(numStripes, ledsPerStripe int) *UdpStripeLamp {
	stripes := make([]Stripe, numStripes)
	for i := range stripes {
		stripes[i] = make(Stripe, ledsPerStripe)
	}
	return &UdpStripeLamp{stripes, nil, make([]uint8, ledsPerStripe*numStripes*3+2), 3}
}

func (l *UdpStripeLamp) Power(on bool) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = 'P'
	if on {
		l.buf[1] = 1

	} else {
		l.buf[1] = 0
	}

	err := l.trans.SendReliable(l.buf[:2])
	return err
}

func (l *UdpStripeLamp) SetBrightness(b uint8) error {
	color := color.RGBA{b, b, b, 0}
	return l.SetColor(&color)
}

func (l *UdpStripeLamp) SetColor(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = 'C'
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf[1], l.buf[2], l.buf[3] = c.R, c.G, c.B
	err := l.trans.SendReliable(l.buf[:4])
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
	err := l.trans.Close()
	l.trans = nil
	return err
}

func (l *UdpStripeLamp) Dial(laddr, raddr *net.UDPAddr) (err error) {

	trans, err := DialReliableUDPTransport(laddr, raddr)
	if err == nil {
		l.trans = trans
	}
	return
}

func (l *UdpStripeLamp) UpdateAll() error {
	if l.trans == nil {
		return errors.New("Not Dialed")

	}
	l.buf[0] = 'D'
	bufpos := 1
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

	err := l.trans.SendReliable(l.buf)
	return err
}
