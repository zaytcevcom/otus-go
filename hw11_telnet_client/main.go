package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Invalid count of arguments")
		return
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to establish connection:", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Connected to %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	messages := make(chan error, 1)

	go func() {
		defer cancel()

		if err := client.Send(); err != nil {
			messages <- err
		}
	}()

	go func() {
		defer cancel()

		if err := client.Receive(); err != nil {
			messages <- err
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "Connection was closed")

		if err := client.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to close connection:", err)
		}
	case message := <-messages:
		fmt.Fprintln(os.Stderr, message)
	}
}
