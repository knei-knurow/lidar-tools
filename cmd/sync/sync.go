package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

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
	flag.IntVar(&baudRate, "baud", 19200, "port baud rate (bps)")
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
	accel := Accel{
		calibration: accelCalib,
		port:        port,
	}
	servo := Servo{
		data:       ServoData{positon: 1600},
		positonMin: 1600,
		positonMax: 3800,
		vector:     25,
		port:       port,
		delayMs:    20,
	}
	lidar := Lidar{ // TODO: make it more configurable from command line
		RPM:  660,
		Mode: rplidarModeDefault,
		Args: "-r 660 -m 2",
		Path: "scan-dummy.exe",
	}

	// Create communication channels
	// lidarChan := make(chan LidarCloud)
	servoChan := make(chan ServoData)
	// accelChan := make(chan AccelData)

	// Start goroutines
	go lidar.StartLoop()
	go servo.StartLoop(servoChan)
	go accel.StartLoop()

	// Main loop
	for {
		select {
		case servoData := <-servoChan:
			if servoOut {
				writer.WriteString(fmt.Sprintf("S %d %d\n", servoData.timept.UnixNano(), servoData.positon))
			}
			// servoBuffer.append(servoData)
			//case accelData := <-accelChan:
			// accelBuffer.append(accelData)
			//case lidarData := <-lidarChan:
			// lidarBuffer.append(lidarData)
		}
		writer.Flush()
	}

	for {

		if accelOut {
			writer.WriteString(fmt.Sprintf("A %d\t%d\t%d\t%d\t%d\t%d\t%d\n", accel.data.timept.UnixNano(),
				accel.data.xAccel, accel.data.yAccel, accel.data.zAccel,
				accel.data.xGyro, accel.data.yGyro, accel.data.zGyro))
		}
	}

	/*
		go Servo.StartLoop(servoChan)
		go Accel.StartLoop(accelChan)
		go Lidar.StartLoop(lidarChan)

		for {
			select {
			case servoData := <-servoChan:
				servoBuffer.append(servoData)
			case accelData := <-accelChan:
				accelBuffer.append(accelData)
			case lidarData := <-lidarChan:
				lidarBuffer.append(lidarData)
			}
		}
	*/
}
