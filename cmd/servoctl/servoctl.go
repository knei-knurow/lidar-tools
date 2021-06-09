package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
	"github.com/knei-knurow/lidar-tools/frames"
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
	flag.UintVar(&baudRate, "baud", 19200, "port baud rate (bps)")
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
		log.Printf("failed to open port %s: %v\n", portName, err)
		return
	}
	defer port.Close()

	if value == -1 {
		fmt.Println("stop because -1 entered")
		return
	}

	if value < int(minValue) {
		log.Printf("warning: %d is smaller than %d\n", value, minValue)
	}

	if value > int(maxValue) {
		log.Printf("warning: %d is bigger than max value %d\n", value, maxValue)
	}

	if value > 65535 {
		log.Printf("error: %d overflows uint16\n", value)
		return
	}

	inputByte := uint16(value)
	data := []byte{byte(inputByte >> 8), byte(inputByte)}
	frame := frames.Create([]byte(frames.LidarHeader), data)

	log.Printf("frame: %s\n", frame)
	for i, currentByte := range frame {
		log.Printf("%d %s will be sent\n", i, frames.DescribeByte(currentByte))
		_, err := port.Write([]byte{currentByte})
		if err != nil {
			log.Printf("%d byte: failed to write it to port: %v\n", i, err)
		}
	}
}
