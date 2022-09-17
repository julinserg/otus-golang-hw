package main

import (
	"bufio"
	"io"
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
	address               string
	timeout               time.Duration
	in                    io.ReadCloser
	out                   io.Writer
	conn                  net.Conn
	inputChNet            chan string
	errorChNet            chan error
	inputChStdIn          chan string
	errorChStdIn          chan error
	sendRoutineIsStart    bool
	receiveRoutineIsStart bool
}

func reader(r io.Reader, inputCh chan string, errorCh chan error) {
	rb := bufio.NewReader(r)
	for {
		dataStr, err := rb.ReadString('\n')
		if err != nil {
			errorCh <- err
			return
		}
		if len(dataStr) == 0 {
			inputCh <- ""
			return
		}
		inputCh <- dataStr
	}
}

func waitReadAndSend(isStartReader *bool, r io.Reader, w io.Writer, inputCh chan string, errorCh chan error) error {
	if !*isStartReader {
		*isStartReader = true
		go reader(r, inputCh, errorCh)
	}

	for {
		select {
		case data := <-inputCh:
			_, err := w.Write([]byte(data))
			return err
		case e := <-errorCh:
			return e
		case <-time.After(1 * time.Second):
			return nil
		}
	}
}

func (tc *TelnetClientImpl) Connect() error {
	var err error
	tc.conn, err = net.DialTimeout("tcp", tc.address, tc.timeout)
	return err
}

func (tc *TelnetClientImpl) Send() error {
	return waitReadAndSend(&tc.sendRoutineIsStart, tc.in, tc.conn, tc.inputChStdIn, tc.errorChStdIn)
}

func (tc *TelnetClientImpl) Receive() error {
	return waitReadAndSend(&tc.receiveRoutineIsStart, tc.conn, tc.out, tc.inputChNet, tc.errorChNet)
}

func (tc *TelnetClientImpl) Close() error {
	return tc.conn.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		address: address, timeout: timeout,
		in: in, out: out,
		inputChNet:   make(chan string),
		errorChNet:   make(chan error),
		inputChStdIn: make(chan string),
		errorChStdIn: make(chan error),
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
