package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	filePath string
	host     string
)

func init() {
	flag.StringVar(&filePath, "file-path", "", ".txt file containing valid cloud series")
	flag.StringVar(&host, "host", "", "ip address of the host which should receive the data")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("sender: failed to open input file")
	}

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
		line = strings.TrimSuffix(line, "\n")

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			go send(host, chunk)
			chunk = make([]byte, 0, 65536)
		} else {
			chunk = append(chunk, []byte(line)...)
		}
	}
}

// Send sends single data to host.
func send(host string, data []byte) {
	fmt.Printf("sender: data chunk of size %d KB\n", len(data)/1024)
}
