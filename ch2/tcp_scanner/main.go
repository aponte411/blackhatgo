package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	// Setup ports and results channels
	ports := make(chan int, 100)
	results := make(chan int)
	// Save open ports to print later
	var openports []int
	// Set up workers for each port connection
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	// send a ports\ to the workers in a seperate goroutine
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// Result gathering loop
	for i := 0; i < 1024; i++ {
		port := <-results
		// If the port is open, add it to the openports slice
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
