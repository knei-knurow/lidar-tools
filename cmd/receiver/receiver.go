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
	log.SetFlags(0)
	log.SetPrefix("receiver: ")

	flag.StringVar(&port, "port", ":8080", "port to listen on")
	flag.BoolVar(&verbose, "verbose", false, "log stuff")
}

func main() {
	flag.Parse()

	pckt, err := net.ListenPacket("udp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v\n", port, err)
	}
	fmt.Fprintf(os.Stderr, "listening on port %s\n", port)
	defer pckt.Close()

	for {
		buf := make([]byte, 65536)
		n, _, err := pckt.ReadFrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read from buffer: %v\n", err)
			break
		}

		text := string(buf[0:n])
		text = strings.TrimSpace(text)
		fmt.Print(text)
	}

	fmt.Fprintln(os.Stderr, "done")
}
