package main

import (
	"context"
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

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &telnetClient{
		address: address,
		in:      in,
		out:     out,
		ctx:     ctx,
		cancel:  cancel,
	}
}

type telnetClient struct {
	address string
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	ctx     context.Context
	cancel  context.CancelFunc
}

func (c *telnetClient) Connect() error {
	var err error
	c.conn, err = (&net.Dialer{Timeout: timeout}).DialContext(c.ctx, "tcp", c.address)
	return err
}

func (c *telnetClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	return err
}
