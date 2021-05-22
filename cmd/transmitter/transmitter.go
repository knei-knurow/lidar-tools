package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	dest string
	port string
)

func init() {
	flag.StringVar(&dest, "dest", "192.168.1.1", "address to send packets to")
	flag.StringVar(&port, "port", "8080", "port on dest to route packets to")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	conn, err := net.Dial("udp", dest+":"+port)
	if err != nil {
		log.Fatalln("transmitter: failed to dial:", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// chunk represents a single cloud scanned from lidar
	chunk := make([]byte, 0, 65536)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("transmitter: end of file")
				break
			}
			log.Fatalln("transmitter:", err)
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			chunk = append(chunk, []byte(line)...)

			cloudIndex, elapsed, err := getCloudData(line)
			if err != nil {
				log.Fatalln("transmitter: failed to get cloud data for line =", line)
			}
			time.Sleep(time.Duration(elapsed) * time.Millisecond)

			go send(conn, chunk, cloudIndex, elapsed)
			chunk = make([]byte, 0, 65536)
			continue
		}

		chunk = append(chunk, []byte(line)...)
	}
}

// Send sends single data to host.
func send(conn net.Conn, data []byte, cloudIndex int, elapsed int) {
	n, err := conn.Write(data)
	if err != nil {
		log.Fatalln("transmitter: failed to write to conn:", err)
	}

	fmt.Printf("transmitter: sent chunk of size %d KB (cloud %d, t %d)\n", n/1024, cloudIndex, elapsed)
}
