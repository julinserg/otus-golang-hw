package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	cancel  context.CancelFunc
}

func (tc *TelnetClientImpl) Connect() error {
	var err error
	tc.conn, err = net.DialTimeout("tcp", tc.address, tc.timeout)
	return err
}

func (tc *TelnetClientImpl) Send() error {
	buffer := make([]byte, 1024)
	numBytes, err := tc.in.Read(buffer)
	if err != nil {
		return err
	}
	if numBytes <= 0 {
		return nil
	}

	numBytes, errNet := tc.conn.Write(buffer[:numBytes])
	if numBytes > 0 {
		log.Printf("To server %v\n", buffer[:numBytes])
	}
	return errNet
}

func (tc *TelnetClientImpl) Receive() error {
	buffer := make([]byte, 1024)
	numBytes, err := tc.conn.Read(buffer)
	if err != nil {
		return err
	}
	if numBytes <= 0 {
		return nil
	}
	log.Printf("From server %v\n", buffer[:numBytes])
	_, err = tc.out.Write(buffer[:numBytes])
	return err
}

func (tc *TelnetClientImpl) Close() error {
	return tc.conn.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{address: address, timeout: timeout, in: in, out: out}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
