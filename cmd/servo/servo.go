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
	portName   string
	baudRate   uint
	dataBits   uint
	stopBits   uint
	parityMode int
)

var (
	value           int
	minValue        uint
	maxValue        uint
	waitForResponse bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&portName, "port", "/dev/tty.*", "port to listen on")
	flag.UintVar(&baudRate, "baud", 9600, "baud rate in bits per second")
	flag.UintVar(&dataBits, "dbits", 8, "the number of data bits in a single frame")
	flag.UintVar(&stopBits, "sbits", 1, "the number of stop bits in a single frame")
	flag.IntVar(&parityMode, "pmode", 1, "parity mode, none = 0, odd = 1, even = 2")
	flag.IntVar(&value, "value", -1, "value to encode into a frame and send")
	flag.UintVar(&minValue, "min-value", 1600, "minimum value that is valid (uint16")
	flag.UintVar(&maxValue, "max-value", 4400, "maximum value that is valid (uint16")
	flag.BoolVar(&waitForResponse, "wait", false, "whether to wait for the MCU to respond with a single byte")
}

func main() {
	flag.Parse()

	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        dataBits,
		StopBits:        stopBits,
		MinimumReadSize: 1,
		ParityMode:      serial.ParityMode(parityMode),
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("servo: failed to open serial port %s: %v\n", portName, err)
	}
	defer port.Close()

	if value == -1 {
		fmt.Println("servo: finish: -1 entered")
		os.Exit(0)
	}

	if value < int(minValue) {
		log.Fatalf("servo: error: %d is smaller than %d\n", value, minValue)
	}

	if value > int(maxValue) {
		log.Fatalf("servo: error: %d is bigger than max value %d\n", value, maxValue)
	}

	if value > 65535 {
		log.Fatalf("servo: error: %d overflows uint16\n", value)
	}

	inputByte := uint16(value)
	f := frame.EncodeFrame(inputByte)

	fmt.Printf("servo: frame: %s\n", f)
	for i, currentByte := range f {
		fmt.Println("---")
		fmt.Printf("servo: %d %s will be sent\n", i, frame.DescribeByte(currentByte))
		_, err := port.Write([]byte{currentByte})
		if err != nil {
			log.Fatalf("servo: %d byte: failed to write it to serial port: %v\n", i, err)
		}
		fmt.Printf("servo: %d byte: wrote it to serial port\n", i)
	}

	if waitForResponse {
		fmt.Printf("servo: waiting for single byte...\n")
		output := make([]byte, 1)
		n, err := port.Read(output)
		if err != nil {
			log.Fatalln("servo: failed to read from serial port:", err)
		}
		outputByte := output[0]

		fmt.Printf("servo: read %d bytes (\"%d\") from serial port \n", n, outputByte)
	}

	fmt.Println("servo: finish")
}
