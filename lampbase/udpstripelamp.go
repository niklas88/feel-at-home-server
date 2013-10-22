package lampbase

import (
	"bytes"
	"errors"
	"image/color"
	"net"
	"time"
)

const maxTries = 4

type UdpStripeLamp struct {
	stripes []Stripe
	raddr   *net.UDPAddr
	laddr   *net.UDPAddr
	conn    *net.UDPConn
	buf     []byte
	seqNum  uint8
}

// Bufs 0th value will be overwritten with the sequence number
func (l *UdpStripeLamp) sendReliable(buf []uint8) error {
	var (
		err    error
		read   int
		ackBuf [4]byte
	)
	// TODO: Log all errors
	success := false
	tries := 0
	l.seqNum++
	l.buf[0] = l.seqNum
	for !success && tries <= maxTries {
		tries++
		_, err = l.conn.Write(buf)
		if err == nil {
			// Try waiting for ACK
			l.conn.SetReadDeadline(time.Now().Add(600 * time.Millisecond))
			for !success {
				read, err = l.conn.Read(ackBuf[:])
				if err != nil {
					if err.(*net.OpError).Timeout() {
						err = errors.New("No ack received")
					}
					break
				}

				if read != 4 || bytes.Equal(ackBuf[:3], []byte("ACK")) {
					err = errors.New("Ack broken: " + string(ackBuf[:]))
				} else {
					// We just ignore/drop non matching ACKs they are old
					success = ackBuf[3] == l.seqNum
				}
			}
		}

	}
	return err
}

func NewUdpStripeLamp(numStripes, ledsPerStripe int) *UdpStripeLamp {
	stripes := make([]Stripe, numStripes)
	for i := range stripes {
		stripes[i] = make(Stripe, ledsPerStripe)
	}
	return &UdpStripeLamp{stripes, nil, nil, nil, make([]uint8, ledsPerStripe*numStripes*3+2), 3}
}

func (l *UdpStripeLamp) Power(on bool) error {
	if l.conn == nil {
		return errors.New("Not Dialed")
	}
	l.buf[1] = 'P'
	if on {
		l.buf[2] = 1

	} else {
		l.buf[2] = 0
	}

	err := l.sendReliable(l.buf[:3])
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
	l.buf[1] = 'C'
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf[2], l.buf[3], l.buf[4] = c.R, c.G, c.B
	err := l.sendReliable(l.buf[:5])
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

func (l *UdpStripeLamp) Dial(laddr, raddr *net.UDPAddr) (err error) {
	l.raddr, l.laddr = raddr, laddr

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

	err := l.sendReliable(l.buf)
	return err
}
