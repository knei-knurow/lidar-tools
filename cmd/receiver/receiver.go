package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for true {
		str, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Fatalln("receiver: data stream ended")
			} else {
				log.Fatalln("receiver: unknown error:", err)
			}
		}

		str = strings.TrimSuffix(str, "\n")

		fmt.Printf("receiver: new data: %s\n", str)
	}
}
