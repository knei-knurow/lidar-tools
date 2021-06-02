package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/knei-knurow/lidar-tools/frames"
	"github.com/tarm/serial"
)

var (
	portName string
	baudRate int
	accelOut bool
	servoOut bool
	lidarOut bool
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("sync: ")

	flag.StringVar(&portName, "port", "COM9", "serial communication port")
	flag.IntVar(&baudRate, "baud", 9600, "port baud rate (bps)")
	flag.BoolVar(&accelOut, "accel", true, "print accelerometer data on stdout")
	flag.BoolVar(&servoOut, "servo", true, "print set servo position on stdout")
	flag.BoolVar(&lidarOut, "lidar", true, "print lidar data on stdout")

	flag.Parse()
	log.Println("starting...")
}

func main() {
	writer := bufio.NewWriter(os.Stdout)

	config := &serial.Config{
		Name: portName,
		Baud: baudRate,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Println("cannot open port:", err)
		return
	}
	defer port.Close()

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
		frame := make(frames.Frame, 18)
		if err := readAceelFrame(port, frame, 'L'); err != nil {
			log.Printf("error: %s\n", err)
		}

		// Accelerometer: Processing data
		accel, err = processAccelFrame(frame)
		if err != nil {
			log.Println("cannot process frame")
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
