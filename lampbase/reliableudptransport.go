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
	addr   *net.UDPAddr
	seqNum uint8
	buf    bytes.Buffer
}

func (l *ReliableUDPTransport) Write(b []byte) (int, error) {
	var (
		ackBuf  [4]byte
		written int
		err     error
	)
	tries := 0
	l.seqNum++
	l.buf.Reset()
	// Note that bytes.Buffer's Write() always returns nil errors
	l.buf.WriteByte(byte(l.seqNum))
	l.buf.Write(b)

	for tries <= maxTries {
		tries++
		written, err = l.conn.WriteToUDP(l.buf.Bytes(), l.addr)
		if written != len(b)+1 {
			return written, fmt.Errorf("could not send as single packet")
		}
		written -= 1 // the seqNum is not part of the data written by the user

		if err == nil {
			// Try waiting for ACK
			l.conn.SetReadDeadline(time.Now().Add(600 * time.Millisecond))
			read, addr, err := l.conn.ReadFrom(ackBuf[:])
			if err != nil {
				err = fmt.Errorf("no ack received %q from %q err: %s", ackBuf, addr, err)
			}
			if read != 4 || !bytes.Equal(ackBuf[:3], []byte("ACK")) {
				err = fmt.Errorf("Ack broken: %q", ackBuf[:])
			} else 	if ackBuf[3] == l.seqNum { // We ignore non matching acks and are done for matching ones
				err = nil
				break;
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
	l.conn, err = net.ListenUDP("udp4", laddr)
	if err == nil {
		l.addr = raddr
	}
	return
}
