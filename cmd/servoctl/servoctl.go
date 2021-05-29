package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jacobsa/go-serial/serial"
	"github.com/knei-knurow/lidar-tools/frame"
)

var (
	portName string
	baudRate uint
)

var (
	value           int
	minValue        uint
	maxValue        uint
	waitForResponse bool
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("servoctl: ")

	flag.StringVar(&portName, "port", "/dev/ttyUSB0", "serial communication port")
	flag.UintVar(&baudRate, "baud", 9600, "port baud rate (bps)")
	flag.IntVar(&value, "value", -1, "value to encode into a frame and send")
	flag.UintVar(&minValue, "min-value", 1600, "minimum value that is valid (uint16)")
	flag.UintVar(&maxValue, "max-value", 4400, "maximum value that is valid (uint16)")
	flag.BoolVar(&waitForResponse, "wait", false, "wait for MCU to respond with a single byte")
}

func main() {
	flag.Parse()

	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
		ParityMode:      0, // no parity
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("failed to open port %s: %v\n", portName, err)
	}
	defer port.Close()

	if value == -1 {
		fmt.Println("finish: -1 entered")
		os.Exit(0)
	}

	if value < int(minValue) {
		log.Fatalf("error: %d is smaller than %d\n", value, minValue)
	}

	if value > int(maxValue) {
		log.Fatalf("error: %d is bigger than max value %d\n", value, maxValue)
	}

	if value > 65535 {
		log.Fatalf("error: %d overflows uint16\n", value)
	}

	inputByte := uint16(value)
	f := frame.EncodeFrame(inputByte)

	fmt.Printf("frame: %s\n", f)
	for i, currentByte := range f {
		fmt.Println("---")
		fmt.Printf("%d %s will be sent\n", i, frame.DescribeByte(currentByte))
		_, err := port.Write([]byte{currentByte})
		if err != nil {
			log.Fatalf("%d byte: failed to write it to port: %v\n", i, err)
		}
		fmt.Printf("%d byte: wrote it to port\n", i)
	}

	// FIXME: doesn't work – output always contains only zeros
	if waitForResponse {
		fmt.Printf("waiting for 2 bytes...\n")
		output := make([]byte, 2)
		n, err := port.Read(output)
		if err != nil {
			log.Fatalln("failed to read from port:", err)
		}
		fmt.Printf("read %d bytes from port\n", n)

		var fullValue uint16
		fullValue = uint16(output[0]) << 8
		fullValue += uint16(output[1])

		for i, b := range output {
			fmt.Printf("%d %s\n", i, frame.DescribeByte(b))
		}

		fmt.Printf("full value (uint16): %d\n", fullValue)
	}

	fmt.Println("finish")
}
