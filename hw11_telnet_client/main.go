package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func readRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN")
				break OUTER
			}
			text := scanner.Text()
			log.Printf("From server: %s", text)
		}
	}
	log.Printf("Finished readRoutine")
}

func writeRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup, stdin chan string) {
	defer wg.Done()
	//scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		case str := <-stdin:
			//if !scanner.Scan() {
			//	break OUTER
			//}
			//str := scanner.Text()
			log.Printf("To server %v\n", str)

			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
	log.Printf("Finished writeRoutine")
}

func stdinScan() chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if scanner.Err() != nil {
			close(out)
		}
	}()
	return out
}

func main() {

	client := NewTelnetClient("127.0.0.1:4242", 10*time.Second, os.Stdin, os.Stdout)
	defer client.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	err := client.Connect()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	} else {
		log.Println("Connected")
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			err := client.Send()
			if err != nil {
				log.Fatalf("Error send: %v", err)
			}
			time.Sleep(10 * time.Microsecond)
		}
	}()

	wg.Add(1)
	go func() {
		for {
			client.Receive()
			if err != nil {
				log.Fatalf("Error receive: %v", err)
			}
			time.Sleep(10 * time.Microsecond)
		}
	}()

	/*ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println(ctx.Err()) // prints "context canceled"
				client.Close()
				log.Fatalln("Ctrl+D press")
				stop()
			}
		}
	}()*/

	wg.Wait()
}
