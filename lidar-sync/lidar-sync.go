package main

import (
	"flag"
	"log"

	"github.com/jacobsa/go-serial/serial"
	"github.com/knei-knurow/lidar-tools/frame"
)

var (
	port     string
	baudrate uint
	verbose  bool
)

func flags() {
	flag.StringVar(&port, "port", "COM9", "Communication port")
	flag.UintVar(&baudrate, "baudrate", 9600, "Communication port baudrate")
	flag.BoolVar(&verbose, "verbose", true, "Be more verbose.")
	flag.Parse()
}

func main() {
	flags()

	log.SetPrefix("lidar-sync: ")
	log.Println("Starting...")

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

	running := true
	for running {
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
		accel, err := processAccelFrame(&frm)
		if err != nil {
			continue
		}
		if verbose {
			log.Printf("accel: %d %d %d %d %d %d\n",
				accel.xAccel, accel.yAccel, accel.zAccel,
				accel.xGyro, accel.yGyro, accel.zGyro)
		}

	}

}
