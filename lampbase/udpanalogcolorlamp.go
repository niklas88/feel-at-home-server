package lampbase

import (
	"errors"
	"image/color"
	"net"
)

type UdpAnalogColorLamp struct {
	trans      *ReliableUDPTransport
	devicePort uint8
	buf        []uint8
}

func NewUdpAnalogColorLamp() *UdpAnalogColorLamp {
	return &UdpAnalogColorLamp{nil, 0, make([]uint8, 5)}
}

func (l *UdpAnalogColorLamp) Power(on bool) error {
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
	_, err := l.trans.Write(l.buf[:3])
	return err
}

func (l *UdpAnalogColorLamp) SetBrightness(b uint8) error {
	color := color.RGBA{b, b, b, 0}
	return l.SetColor(&color)
}

func (l *UdpAnalogColorLamp) SetColor(col color.Color) error {
	if l.trans == nil {
		return errors.New("Not Dialed")
	}
	l.buf[0] = l.devicePort
	l.buf[1] = 'C'
	c := color.RGBAModel.Convert(col).(color.RGBA)
	l.buf[2], l.buf[3], l.buf[4] = c.R, c.G, c.B
	_, err := l.trans.Write(l.buf[:5])
	return err
}

func (l *UdpAnalogColorLamp) Close() error {
	err := l.trans.Close()
	l.trans = nil
	return err
}

func (l *UdpAnalogColorLamp) Dial(laddr, raddr *net.UDPAddr, devicePort uint8) (err error) {
	l.devicePort = devicePort

	trans, err := DialReliableUDPTransport(laddr, raddr)
	if err == nil {
		l.trans = trans
	}
	return
}
