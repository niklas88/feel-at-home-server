package lampbase

import (
	"errors"
	"image/color"
	"net"
	"time"
)

type UdpWordClock struct {
	UdpMatrixLamp
	timeUpdateInterval time.Duration
}

func NewUdpWordClock(timeUpdateInterval time.Duration) *UdpWordClock {
	lamp := new(UdpWordClock)
	return lamp
}

func (l *UdpWordClock) Dial(laddr, raddr *net.UDPAddr, lampNum uint8) error {
	err := l.UdpPowerDevice.Dial(laddr, raddr, lampNum)
	if err == nil {
		go func() {
			c := time.Tick(l.timeUpdateInterval)
			for now := range c {
				if l.trans != nil {
					l.TimeUpdate(now)
				}
			}
		}()
	}
	return err
}

func (l *UdpMatrixLamp) Clock() error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('T', 0x00)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpWordClock) ClockColor(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpDimLamp.writeHead('T', 0x01)
	l.writeColor(col)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpWordClock) TimeUpdate(t time.Time) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.UdpColorLamp.writeHead('T', 0x02)
	l.writeTime(t)
	_, err := l.buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) writeTime(t time.Time) {
	unix := t.Unix()
	l.buf.WriteByte(byte(unix >> 56))
	l.buf.WriteByte(byte((unix >> 48) & 0xff))
	l.buf.WriteByte(byte((unix >> 40) & 0xff))
	l.buf.WriteByte(byte((unix >> 32) & 0xff))
	l.buf.WriteByte(byte((unix >> 24) & 0xff))
	l.buf.WriteByte(byte((unix >> 16) & 0xff))
	l.buf.WriteByte(byte((unix >> 8) & 0xff))
	l.buf.WriteByte(byte(unix & 0xff))
}
