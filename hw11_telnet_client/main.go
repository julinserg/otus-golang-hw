package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

func main() {

	client := NewTelnetClient("192.168.1.145:4242", 5*time.Second, os.Stdin, os.Stdout)
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())

	err := client.Connect()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	} else {
		log.Println("Connected")
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Send done")
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
				log.Println("Receive done")
				return
			default:
				client.Receive()
				log.Println("Receive")
				if err != nil {
					log.Printf("Error receive: %v \n", err)
					cancel()
					return
				}
				time.Sleep(10 * time.Microsecond)
			}
		}
	}()

	/*ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				cancel()
				log.Fatalln("Correct exit")
				stop()
			}
		}
	}()*/

	wg.Wait()
	log.Println("TelnetClient exit")
}
