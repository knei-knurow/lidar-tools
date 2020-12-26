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
	filePath string
	port     string
	dest     string
)

func init() {
	flag.StringVar(&filePath, "file-path", "", ".txt file containing valid cloud series")
	flag.StringVar(&dest, "dest", "192.168.1.1", "address to send packets to")
	flag.StringVar(&port, "port", ":8080", "port on destAddress to route packets to")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("sender: failed to open input file")
	}

	conn, err := net.Dial("udp", dest+port)
	if err != nil {
		log.Fatalln("sender: failed to dial")
	}
	conn.SetDeadline(time.Now().Add(time.Second)) // too long anyway
	defer conn.Close()

	reader := bufio.NewReader(f)

	// chunk represents a single cloud scanned from lidar
	chunk := make([]byte, 0, 65536)
	for true {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("sender: end of file")
				break
			}
			log.Fatalln("sender:", err)
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			go send(conn, chunk)
			chunk = make([]byte, 0, 65536)
		} else {
			chunk = append(chunk, []byte(line)...)
		}
	}
}

// Send sends single data to host.
func send(conn net.Conn, data []byte) {
	n, err := conn.Write(data)
	if err != nil {
		log.Fatalln("sender: failed to write to conn")
	}

	fmt.Printf("sender: data chunk of size %d KB\n", n/1024)
}
