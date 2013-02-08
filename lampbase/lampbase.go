package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type Stripe []color.RGBA

type Lamp struct {
	Stripes []Stripe
	Addr    *net.UDPAddr
	conn    *net.UDPConn
	buf     []uint8
}

func NewLamp(numStripes, ledsPerStripe int, addr *net.UDPAddr) *Lamp {
	stripes := make([]Stripe, numStripes)
	for i, _ := range stripes {
		stripes[i] = make(Stripe, ledsPerStripe)
	}
	return &Lamp{stripes, addr, nil, make([]uint8, ledsPerStripe*numStripes*3+1)}
}

func (l *Lamp) Close() {
	l.conn.Close()
	l.conn = nil
}

func (l *Lamp) Dial() (err error) {
	conn := l.conn

	if conn == nil {
		conn, err := net.DialUDP("udp4", nil, l.Addr)
		if err == nil {
			l.conn = conn
		}
	}
	return
}

func (l *Lamp) Update() (err error) {
	if l.conn == nil {
		if err = l.Dial(); err != nil {
			return
		}
	}
	l.buf[0] = 'D'
	bufpos := 1
	for i, s := range l.Stripes {
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
	return
}
