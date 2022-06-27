package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	_host       string
	_nb_workers int
	_channels   chan int
)

func init() {
	flag.StringVar(&_host, "h", "127.0.0.1", "machine is the name of the machine")
	flag.IntVar(&_nb_workers, "w", 10, "number of workers")
	flag.Parse()
}

// Take a port number as argument and try to connect to it with the hostname given in the -h flag.
func connectTCP(port int) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", _host, port))
	if err != nil {
		return
	}
	fmt.Println("Connection successful in port ", port)
	if err = c.Close(); err != nil {
		fmt.Println("Error closing connection: ", err)
		return
	}
}

func worker() {
	for p := range _channels {
		connectTCP(p)
	}
}

func main() {
	// Create a channel to communicate with the workers
	if _nb_workers <= 0 {
		_nb_workers = 10
	}
	_channels = make(chan int, _nb_workers)

	// Start the workers
	for i := 0; i < cap(_channels); i++ {
		go worker()
	}

	// Send the ports to the workers
	for i := 1; i <= 3000; i++ {
		_channels <- i
	}

	close(_channels)
}
