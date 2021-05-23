package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bartekpacia/lidar-tools/frame"
	"github.com/jacobsa/go-serial/serial"
)

var (
	portName   string
	baudRate   uint
	dataBits   uint
	stopBits   uint
	parityMode int
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&portName, "port", "/dev/tty.*", "port to listen on")
	flag.UintVar(&baudRate, "baud", 9600, "baud rate in bits per second")
	flag.UintVar(&dataBits, "dbits", 8, "the number of data bits in a single frame")
	flag.UintVar(&stopBits, "sbits", 1, "the number of stop bits in a single frame")
	flag.IntVar(&parityMode, "pmode", 1, "parity mode, none = 0, odd = 1, even = 2")
}

func main() {
	flag.Parse()

	// temporarily - just for testing
	// fmt.Printf("\x10 in ascii: %q\n", 0)
	// rf0 := frame.CreateRawFrame(0)
	// f0 := frame.CreateFrame(0)
	// crc0 := frame.CalculateCRC(rf0)
	// fmt.Printf("rawFrame0: % x\n", rf0)
	// fmt.Printf("frame0: % x\n", f0)
	// fmt.Printf("crc0: % x\n", crc0)

	// fmt.Println("---")
	// rf1 := frame.CreateRawFrame(3)
	// f1 := frame.CreateFrame(3)
	// crc1 := frame.CalculateCRC(rf1)
	// fmt.Printf("rawFrame1: % x\n", rf1)
	// fmt.Printf("frame1: % x\n", f1)
	// fmt.Printf("crc1: % x\n", crc1)

	// fmt.Println("---")
	// rf2 := frame.CreateRawFrame(5)
	// f2 := frame.CreateFrame(5)
	// crc2 := frame.CalculateCRC(rf2)
	// fmt.Printf("rawFrame2: % x\n", rf2)
	// fmt.Printf("frame2: % x\n", f2)
	// fmt.Printf("crc2: % x\n", crc2)

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
		log.Fatalf("failed to open serial port %s: %v\n", portName, err)
	}
	defer port.Close()

	fmt.Println("uart_echo: a tiny program to control a servo")
	fmt.Println("uart_echo: enter -1 to stop")

	for {
		var value int
		fmt.Print("enter a single 16-bit number to be sent (integer, 0-65536): ")
		_, err := fmt.Scanf("%d", &value)
		if err != nil {
			log.Fatalln("error reading from stdin:", err)
		}

		if value == -1 {
			break
		}

		if value < 1600 || value > 4400 {
			fmt.Printf("uart_echo: error: %d overflows uint16\n", value)
			break
		}

		inputByte := uint16(value)
		f := frame.CreateFrame(inputByte)

		fmt.Printf("frame: %s\n", f)
		for i, currentByte := range f {
			fmt.Println("---")
			fmt.Printf("%d %s will be sent\n", i, frame.DescribeByte(currentByte))
			_, err := port.Write([]byte{currentByte})
			if err != nil {
				log.Fatalf("%d byte: failed to write it to serial port: %v\n", i, err)
			}
			fmt.Printf("%d byte: wrote it to serial port\n", i)
		}

		// Receiving - commented out
		// output := make([]byte, 1)
		// n, err = port.Read(output)
		// if err != nil {
		// 	log.Fatalln("error reading from serial port:", err)
		// }
		// outputByte := output[0]

		// fmt.Printf("read %d bytes (\"%d\") from serial port \n", n, outputByte)
	}

	fmt.Println("finished :)")
}
