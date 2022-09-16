package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second, "Connection timeout")
}

func main() {
	flag.Parse()
	if len(os.Args) < 3 {
		log.Fatalln("ERROR: Set address and port")
	}
	client := NewTelnetClient(net.JoinHostPort(os.Args[1], os.Args[2]), timeout, os.Stdin, os.Stdout)
	isConnected := false
	defer func(isConn *bool) {
		if *isConn {
			client.Close()
		}
	}(&isConnected)

	ctx, cancel := context.WithCancel(context.Background())

	if err := client.Connect(); err != nil {
		log.Printf("Cannot connect: %v \n", err)
		return
	}
	isConnected = true
	log.Println("Connected")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			log.Println("'Send goroutine' done")
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Send()
				if err != nil {
					log.Printf("Error send: %v \n", err)
					cancel()
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			log.Println("'Receive goroutine' done")
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Receive()
				if err != nil {
					log.Printf("Error receive: %v \n", err)
					// cancel()
					return
				}
			}
		}
	}()

	ctxInterrupt, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	wg.Add(1)
	go func() {
		defer func() {
			log.Println("'Wait Ctrl+C press goroutine' done")
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ctxInterrupt.Done():
				log.Println("Ctrl+C press")
				cancel()
				return
			}
		}
	}()

	wg.Wait()
	log.Println("TelnetClient exit")
}
