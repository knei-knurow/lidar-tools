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
	avrPortName   string
	avrBaudRate   int
	lidarPortName string
	lidarMode     int
	lidarRPM      int
	lidarExe      string
	accelOut      bool
	servoOut      bool
	lidarOut      bool
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("sync: ")

	flag.StringVar(&avrPortName, "avrport", "COM9", "AVR serial communication port")
	flag.IntVar(&avrBaudRate, "avrbaud", 19200, "port baud rate (bps)")

	flag.StringVar(&lidarExe, "lidarexe", "lidar.exe", "lidar-scan executable")
	flag.StringVar(&lidarPortName, "lidarport", "COM4", "RPLIDAR serial communication port")
	flag.IntVar(&lidarMode, "lidarmode", rplidarModeDefault, "RPLIDAR mode")
	flag.IntVar(&lidarRPM, "lidarpm", 660, "RPLIDAR given revolutions per minute")

	flag.BoolVar(&accelOut, "accel", true, "print accelerometer data on stdout")
	flag.BoolVar(&servoOut, "servo", true, "print set servo position on stdout")
	flag.BoolVar(&lidarOut, "lidar", true, "print lidar data on stdout")

	flag.Parse()
	log.Println("starting...")
}

func main() {
	writer := bufio.NewWriter(os.Stdout)

	config := &serial.Config{
		Name: avrPortName,
		Baud: avrBaudRate,
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
		vector:     50,
		port:       port,
		delayMs:    100,
	}
	lidar := Lidar{ // TODO: make it more configurable from command line
		RPM:  lidarRPM,
		Mode: lidarMode,
		Args: []string{lidarPortName, "--rpm", fmt.Sprint(lidarRPM), "--mode", fmt.Sprint(lidarMode)},
		Path: lidarExe, // TODO: Check if exists
	}

	// Create communication channels
	lidarChan := make(chan *LidarCloud) // LidarCloud is >64kB so it cannot be directly passed by a channel
	servoChan := make(chan ServoData)
	accelChan := make(chan AccelData)

	// Start goroutines
	go lidar.StartLoop(lidarChan)
	go servo.StartLoop(servoChan)
	go accel.StartLoop(accelChan)

	// Main loop
	for {
		select {
		case lidarData := <-lidarChan:
			if lidarOut {
				writer.WriteString(fmt.Sprintf("L %d %d\n", lidarData.ID, lidarData.Size))
			}
		case servoData := <-servoChan:
			if servoOut {
				writer.WriteString(fmt.Sprintf("S %d %d\n", servoData.timept.UnixNano(), servoData.positon))
			}
		case accelData := <-accelChan:
			if accelOut {
				writer.WriteString(fmt.Sprintf("A %d\t%d\t%d\t%d\t%d\t%d\t%d\n", accelData.timept.UnixNano(),
					accelData.xAccel, accelData.yAccel, accelData.zAccel,
					accelData.xGyro, accelData.yGyro, accelData.zGyro))
			}
			// accelBuffer.append(accelData)
			//case lidarData := <-lidarChan:
			// lidarBuffer.append(lidarData)
		}
		writer.Flush()
	}
}
