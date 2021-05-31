package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/knei-knurow/lidar-tools/frames"
)

var (
	portName string
	baudRate uint
	accelOut bool
	servoOut bool
	lidarOut bool
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("sync: ")

	flag.StringVar(&portName, "port", "COM9", "serial communication port")
	flag.UintVar(&baudRate, "baud", 9600, "port baud rate (bps)")
	flag.BoolVar(&accelOut, "accel", true, "print accelerometer data on stdout")
	flag.BoolVar(&servoOut, "servo", true, "print set servo position on stdout")
	flag.BoolVar(&lidarOut, "lidar", true, "print lidar data on stdout")

	flag.Parse()
	log.Println("starting...")
}

func main() {
	writer := bufio.NewWriter(os.Stdout)

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
	log.Println("connection established")

	// Sources of data initialization
	accel := AccelData{}
	servo := Servo{positon: 3600, positonMin: 1600, positonMax: 4400, vector: 60}
	lidar := Lidar{ // TODO: make it more configurable from command line
		RPM:  660,
		Mode: rplidarModeDefault,
		Args: "-r 660 -m 2",
		Path: "scan-dummy.exe",
	}

	// Start lidar loop
	go lidar.StartLoop()

	// Start accelerometer/servo loop
	for {
		// Servo: Sending data
		servo.move()
		inputByte := servo.positon
		data := []byte{byte(inputByte >> 8), byte(inputByte)} // TODO: Check whether correct
		f := frames.Create([]byte(frames.LidarHeader), data)
		for i, currentByte := range f {
			if _, err := port.Write([]byte{currentByte}); err != nil {
				log.Fatalf("%d byte: failed to write it to port: %v\n", i, err)
			}
		}
		servo.timept = time.Now()

		// Accelerometer: Reading data
		// TODO: this code is very messy but handles bad data much faster
		ok, end := true, true
		var data strings.Builder
		for ok && end {
			buf := make([]byte, 1)
			_, err = port.Read(buf)
			if err != nil {
				log.Println("failed to read from port:", err)
				continue
			}
			switch {
			case buf[0] == 'L' && data.Len() == 0:
				data.WriteString(string(buf[0]))
			case buf[0] == 'D' && data.Len() == 1:
				data.WriteString(string(buf[0]))
			case buf[0] == '-' && data.Len() == 2:
				data.WriteString(string(buf[0]))
			case buf[0] == '#' && data.Len() == 15:
				data.WriteString(string(buf[0]))
			case buf[0] == 'S' && data.Len() == 16:
				data.WriteString(string(buf[0]))
				end = false
			case data.Len() >= 3 && data.Len() <= 14:
				data.WriteString("x")
			default:
				ok = false
			}
		}
		if !ok {
			if accelOut {
				log.Printf("bad accelerometer data: %d bytes\n", data.Len()+1)
			}
			continue
		}

		// TODO: Clever Buffer -> Frame conversion
		// TODO: Use frames.Assemble
		frame := frames.Frame{
			Header:   frames.FrameHeader(data.String()[0:3]),
			Data:     []byte(data.String()[3:15]),
			Checksum: byte(data.String()[16]),
		}

		// Accelerometer
		accel, err = processAccelFrame(&frame)
		if err != nil {
			// FIXME: Handle error
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
