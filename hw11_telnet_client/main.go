package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	// client := NewTelnetClient("192.168.1.145:4242", 5*time.Second, os.Stdin, os.Stdout)
	client := NewTelnetClient("localhost:4242", 5*time.Second, os.Stdin, os.Stdout)
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
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Send goroutine done")
				return
			default:
				err := client.Send()
				if err != nil {
					log.Printf("Error send: %v \n", err)
					cancel()
					return
				}
				time.Sleep(10 * time.Microsecond)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Receive goroutine done")
				return
			default:
				err := client.Receive()
				if err != nil {
					log.Printf("Error receive: %v \n", err)
					// cancel()
					return
				}
				time.Sleep(10 * time.Microsecond)
			}
		}
	}()

	ctxInterrupt, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Wait Ctrl+C press goroutine done")
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
