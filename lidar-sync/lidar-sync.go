package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jacobsa/go-serial/serial"
	"github.com/knei-knurow/lidar-tools/frame"
)

var (
	port     string
	baudrate uint
)

func flags() {
	flag.StringVar(&port, "port", "COM9", "Communication port")
	flag.UintVar(&baudrate, "baudrate", 9600, "Communication port baudrate")
	flag.Parse()
}

func receiver() {

}

func main() {
	flags()

	log.SetPrefix("lidar-sync: ")
	log.Println("Starting...")

	writer := bufio.NewWriter(os.Stdout)

	// Port
	serOpt := serial.OpenOptions{
		PortName:        port,
		BaudRate:        baudrate,
		DataBits:        8,
		MinimumReadSize: 1,
	}
	ser, err := serial.Open(serOpt)
	defer ser.Close()
	if err != nil {
		log.Println("Unable to open the specified port.")
		return
	}
	log.Println("Connection established.")

	var accel AccelData
	servo := 3600

	// Data reading loop
	for {
		// Sending data
		inputByte := uint16(servo)
		f := frame.EncodeFrame(inputByte)
		for _, currentByte := range f {
			if _, err := ser.Write([]byte{currentByte}); err != nil {
				log.Println("Unable to send servo data:", err)
			}
		}

		// Reading data
		buf := make([]byte, 32)
		_, err = ser.Read(buf)
		if err != nil {
			log.Println("An error occured while reading a buffer:", err)
			continue
		}

		// TODO: Clever Buffer -> Frame conversion
		frm := frame.Frame{
			Header:   frame.FrameHeader(buf[0:3]),
			Data:     buf[3:15],
			Checksum: buf[16],
		}

		// Accelerometer
		accel, err = processAccelFrame(&frm)
		if err != nil {
			continue
		}

		// Write stdout
		writer.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			accel.xAccel, accel.yAccel, accel.zAccel,
			accel.xGyro, accel.yGyro, accel.zGyro, servo))
		writer.Flush()
	}
}
