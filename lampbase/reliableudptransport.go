package lampbase

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

const maxTries = 4

type ReliableUDPTransport struct {
	conn   *net.UDPConn
	seqNum uint8
}

func (l *ReliableUDPTransport) SendReliable(b []uint8) error {
	var (
		err    error
		ackBuf [4]byte
		buf    []byte
	)
	success := false
	tries := 0
	l.seqNum++
	buf = make([]byte, len(b)+1)
	buf[0] = l.seqNum
	copy(buf[1:], b[:])

	for !success && tries <= maxTries {
		tries++
		read, err := l.conn.Write(buf)
		if read != len(buf) {
			return fmt.Errorf("could not send as single packet")
		}

		if err == nil {
			// Try waiting for ACK
			l.conn.SetReadDeadline(time.Now().Add(600 * time.Millisecond))
			for !success {
				read, err := l.conn.Read(ackBuf[:])
				if err != nil && err.(*net.OpError).Timeout() {
					return fmt.Errorf("no ack received %q err: %s", ackBuf, err)
				}

				if read != 4 || !bytes.Equal(ackBuf[:3], []byte("ACK")) {
					return fmt.Errorf("Ack broken: %q", ackBuf[:])
				}

				// We just ignore/drop non matching ACKs they are old
				if ackBuf[3] == l.seqNum {
					success = true
					err = nil
				}
			}
		}

	}
	return err
}

func (l *ReliableUDPTransport) Close() error {
	return l.conn.Close()
}

func DialReliableUDPTransport(laddr, raddr *net.UDPAddr) (l *ReliableUDPTransport, err error) {
	l = new(ReliableUDPTransport)
	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err == nil {
		l.conn = conn
	}
	return
}
