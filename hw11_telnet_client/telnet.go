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
	address            string
	timeout            time.Duration
	in                 io.ReadCloser
	out                io.Writer
	conn               net.Conn
	inputChNet         chan interface{}
	errorChNet         chan error
	inputChStdIn       chan interface{}
	errorChStdIn       chan error
	sendRoutineIsStart bool
}

func reader(r io.Reader, input chan interface{}, errorCh chan error) {
	for {
		buffer := make([]byte, 4100)
		numBytes, err := r.Read(buffer)
		if err != nil {
			errorCh <- err
			return
		}
		if numBytes <= 0 {
			input <- struct{}{}
			return
		}
		input <- buffer[:numBytes]
	}
}

func (tc *TelnetClientImpl) Connect() error {
	var err error
	tc.conn, err = net.DialTimeout("tcp", tc.address, tc.timeout)

	tc.inputChNet = make(chan interface{})
	tc.errorChNet = make(chan error)
	tc.inputChStdIn = make(chan interface{})
	tc.errorChStdIn = make(chan error)

	go reader(tc.conn, tc.inputChNet, tc.errorChNet)

	return err
}

func (tc *TelnetClientImpl) Send() error {
	if !tc.sendRoutineIsStart {
		tc.sendRoutineIsStart = true
		go reader(tc.in, tc.inputChStdIn, tc.errorChStdIn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	for {
		select {
		case data := <-tc.inputChStdIn:
			numBytes, err := tc.conn.Write(data.([]byte))
			if numBytes > 0 {
				log.Printf("To server %v\n", data)
			}
			cancel()
			return err
		case e := <-tc.errorChStdIn:
			cancel()
			return e
		case <-ctx.Done():
			cancel()
			return nil
		}
	}
}

func (tc *TelnetClientImpl) Receive() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	for {
		select {
		case data := <-tc.inputChNet:
			log.Printf("From server %v\n", data)
			_, err := tc.out.Write(data.([]byte))
			cancel()
			return err
		case e := <-tc.errorChNet:
			cancel()
			return e
		case <-ctx.Done():
			cancel()
			return nil
		}
	}
}

func (tc *TelnetClientImpl) Close() error {
	return tc.conn.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{address: address, timeout: timeout, in: in, out: out}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
