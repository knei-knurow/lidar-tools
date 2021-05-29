package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/knei-knurow/lidar-tools/frame"
)

var (
	port     string
	baudrate uint
	accelOut bool
	servoOut bool
	lidarOut bool
)

func init() {
	flag.StringVar(&port, "port", "COM9", "Communication port")
	flag.UintVar(&baudrate, "baudrate", 9600, "Communication port baudrate")
	flag.BoolVar(&accelOut, "accel", true, "Print accelerometer data on stdout")
	flag.BoolVar(&servoOut, "servo", true, "Print set servo position on stdout")
	flag.BoolVar(&lidarOut, "lidar", true, "Print lidar data on stdout")
}

func receiver() {

}

func main() {
	flag.Parse()

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

	accel := AccelData{}
	servo := Servo{positon: 3600, positonMin: 1600, positonMax: 4400, vector: 60}

	// Data reading loop
	for {
		// Servo: Sending data
		servo.move()
		inputByte := uint16(servo.positon)
		f := frame.EncodeFrame(inputByte)
		for _, currentByte := range f {
			if _, err := ser.Write([]byte{currentByte}); err != nil {
				log.Println("Unable to send servo data:", err)
			}
		}
		servo.timept = time.Now()

		// Accelerometer: Reading data
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
		if accelOut {
			writer.WriteString(fmt.Sprintf("A %d\t%d\t%d\t%d\t%d\t%d\t%d\n", accel.timept.UnixNano(),
				accel.xAccel, accel.yAccel, accel.zAccel,
				accel.xGyro, accel.yGyro, accel.zGyro))
		}
		if servoOut {
			writer.WriteString(fmt.Sprintf("S %d %d\n", servo.timept.UnixNano(), servo.positon))
		}
		writer.Flush()
	}
}
