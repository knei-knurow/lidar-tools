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
	"strconv"
	"strings"
	"time"
)

var (
	filePath string
	port     string
	dest     string
)

func init() {
	flag.StringVar(&dest, "dest", "192.168.1.1", "address to send packets to")
	flag.StringVar(&port, "port", "8080", "port on destAddress to route packets to")
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
	for true {
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
			commentLine := strings.Split(line, " ")
			timeInt, _ := strconv.Atoi(commentLine[1])
			timeout := time.Duration(timeInt) * time.Millisecond
			time.Sleep(timeout)

			chunk = make([]byte, 0, 65536)
			go send(conn, chunk)
		} else {
			chunk = append(chunk, []byte(line)...)
		}
	}
}

// Send sends single data to host.
func send(conn net.Conn, data []byte) {
	n, err := conn.Write(data)
	if err != nil {
		log.Fatalln("transmitter: failed to write to conn:", err)
	}

	fmt.Printf("transmitter: sent chunk of size %d KB\n", n/1024)
}
