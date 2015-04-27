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
	buf    bytes.Buffer
}

func (l *ReliableUDPTransport) Write(b []byte) (int, error) {
	var (
		ackBuf [4]byte
	)
	tries := 0
	l.seqNum++
	l.buf.Reset()
	err := l.buf.WriteByte(byte(l.seqNum))
	if err != nil {
		return 0, err // zero because we didn't send anything on the network
	}
	written, err := l.buf.Write(b)
	if written != len(b) || err != nil {
		return 0, err
	}

	for tries <= maxTries {
		tries++
		written, err = l.conn.Write(l.buf.Bytes())
		if written != len(b)+1 {
			return written, fmt.Errorf("could not send as single packet")
		}

		if err == nil {
			// Try waiting for ACK
			l.conn.SetReadDeadline(time.Now().Add(600 * time.Millisecond))
			success := false
			for !success {
				read, err := l.conn.Read(ackBuf[:])
				if err != nil {
					return written, fmt.Errorf("no ack received %q err: %s", ackBuf, err)
				}

				if read != 4 || !bytes.Equal(ackBuf[:3], []byte("ACK")) {
					return written, fmt.Errorf("Ack broken: %q", ackBuf[:])
				}

				// We ignore non matching acks and are done for matching ones
				if ackBuf[3] == l.seqNum {
					success = true
					err = nil
				}
			}
		}

	}
	return written, err
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
