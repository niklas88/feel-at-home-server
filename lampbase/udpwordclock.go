package lampbase

import (
	"bytes"
	"errors"
	"image/color"
	"net"
	"time"
)

type UdpWordClock struct {
	UdpMatrixLamp
	updateInterval time.Duration
}

func NewUdpWordClock(timeUpdateInterval time.Duration) *UdpWordClock {
	return &UdpWordClock{updateInterval: timeUpdateInterval}
}

func (l *UdpWordClock) Dial(laddr, raddr *net.UDPAddr, lampNum uint8) error {
	err := l.UdpPowerDevice.Dial(laddr, raddr, lampNum)
	if err == nil {
		go func() {
			c := time.Tick(l.updateInterval)
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
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('T', 0x00, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpWordClock) ClockColor(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpDimLamp.writeHead('T', 0x01, &buf)
	l.writeColor(col, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpWordClock) TimeUpdate(t time.Time) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	var buf bytes.Buffer
	l.UdpColorLamp.writeHead('T', 0x02, &buf)
	l.writeTime(t, &buf)
	_, err := buf.WriteTo(l.trans)
	return err
}

func (l *UdpDimLamp) writeTime(t time.Time, buf *bytes.Buffer) {
	unix := t.Unix()
	_, offset := t.Zone()
	local := unix + int64(offset)
	buf.WriteByte(byte(local >> 56))
	buf.WriteByte(byte((local >> 48) & 0xff))
	buf.WriteByte(byte((local >> 40) & 0xff))
	buf.WriteByte(byte((local >> 32) & 0xff))
	buf.WriteByte(byte((local >> 24) & 0xff))
	buf.WriteByte(byte((local >> 16) & 0xff))
	buf.WriteByte(byte((local >> 8) & 0xff))
	buf.WriteByte(byte(local & 0xff))
}
