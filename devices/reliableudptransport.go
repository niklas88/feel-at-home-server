package devices

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

const maxTries = 4

// ReliableUDPTransport implements a UDP based message oriented lowlevel
// (currently one-way) protocoll by implementing resend of unackowledged
// messages.
type ReliableUDPTransport struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	seqNum uint8
}

// Write implements the Writer interface
func (l *ReliableUDPTransport) Write(b []byte) (written int, lastError error) {
	var (
		ackBuf [4]byte
		err    error
		buf    bytes.Buffer
	)
	tries := 0
	lastError = nil
	l.seqNum++
	// Note that bytes.Buffer's Write() always returns nil errors
	buf.WriteByte(byte(l.seqNum))
	buf.Write(b)

	for tries <= maxTries {
		tries++
		written, err = l.conn.WriteToUDP(buf.Bytes(), l.addr)
		if err != nil {
			lastError = err
			continue
		}

		if written != len(b)+1 {
			lastError = fmt.Errorf("could not send as single packet")
			continue
		}
		written -= 1 // the seqNum is not part of the data written by the user

		// Try waiting for ACK
		l.conn.SetReadDeadline(time.Now().Add(600 * time.Millisecond))
		read, addr, err := l.conn.ReadFrom(ackBuf[:])
		if err != nil {
			lastError = fmt.Errorf("no ack received %q from %q err: %s", ackBuf, addr, err)
			continue
		}

		if read < 4 || !bytes.Equal(ackBuf[:3], []byte("ACK")) {
			log.Println("Try %d : received  Ack: %q and expected seqNum %q", read, ackBuf[:], l.seqNum)
			lastError = fmt.Errorf("Ack broken: %q", ackBuf[:])
			continue
		}

		// Note that if the seqNum doesn't match we retry
		if ackBuf[3] == l.seqNum {
			lastError = nil
			break
		}

	}
	return written, lastError
}

// Close implements the Closer interface
func (l *ReliableUDPTransport) Close() error {
	return l.conn.Close()
}

// DialReliableUDPTransport allows dialing remote UDPAddr
func DialReliableUDPTransport(laddr, raddr *net.UDPAddr) (l *ReliableUDPTransport, err error) {
	l = new(ReliableUDPTransport)
	l.conn, err = net.ListenUDP("udp4", laddr)
	if err == nil {
		l.addr = raddr
	}
	return
}
