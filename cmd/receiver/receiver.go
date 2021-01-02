package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var port string
var verbose bool

func init() {
	flag.StringVar(&port, "port", ":8080", "port to listen on")
	flag.BoolVar(&verbose, "verbose", false, "whether to log stuff")
}

func main() {
	flag.Parse()

	pckt, err := net.ListenPacket("udp", port)
	if err != nil {
		log.Fatalf("receiver: error listening on port %s: %v\n", port, err)
	}
	fmt.Fprintf(os.Stderr, "receiver: listening on port %s\n", port)
	defer pckt.Close()

	for {
		buf := make([]byte, 65536)
		n, _, err := pckt.ReadFrom(buf)
		if err != nil {
			log.Fatalf("receiver: error reading from buffer: %v\n", err)
			break
		}

		text := string(buf[0:n])
		text = strings.TrimSpace(text)
		fmt.Fprintf(os.Stderr, "receiver: received %d KB of data\n", n/1024)
	}

	fmt.Fprintln(os.Stderr, "receiver: done")
}
