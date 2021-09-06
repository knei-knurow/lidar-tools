package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

var (
	// AVR args
	avrPort     string
	avrBaudRate int

	// Lidar flags
	lidarPort string
	lidarMode int
	lidarRPM  int
	lidarExe  string

	// Misc args
	cloudRotation float64
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("sync: ")

	// AVR args
	flag.StringVar(&avrPort, "avrport", "COM13", "AVR serial communication port")
	flag.IntVar(&avrBaudRate, "avrbaud", 19200, "port baud rate (bps)")

	// Lidar args
	flag.StringVar(&lidarExe, "lidarexe", "lidar.exe", "lidar-scan executable")
	flag.StringVar(&lidarPort, "lidarport", "COM4", "RPLIDAR serial communication port")
	flag.IntVar(&lidarMode, "lidarmode", rplidarModeDefault, "RPLIDAR mode (3 - best for indoor, 4 - best for outdoor)")
	flag.IntVar(&lidarRPM, "lidarpm", 660, "RPLIDAR given revolutions per minute")

	// Misc args
	flag.Float64Var(&cloudRotation, "cloudrotation", PrototypeCloudRotation, "each scanned 2D cloud will be rotated by CloudRotation radians, this value depends on the physical lidar location")

	flag.Parse()
	log.Println("starting...")
}

func main() {
	writer := bufio.NewWriter(os.Stdout)

	config := &serial.Config{
		Name: avrPort,
		Baud: avrBaudRate,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Println("cannot open AVR port:", err)
		return
	}
	defer port.Close()

	// Sources of data initialization
	accel := Accel{
		calibration: accelCalib,
		accelScale:  AccelScaleDefault,
		gyroScale:   GyroScaleDefault,
		deltaTime:   DeltaTimeDefault,
		port:        port,
		mode:        AccelModeRaw,
	}
	servo := Servo{
		data:       ServoData{positon: servoStartPos},
		positonMin: servoMinPos,
		positonMax: servoMaxPos,
		vector:     2,
		port:       port,
		delayMs:    40,
	}
	log.Println("servo is setting to the start position")
	servo.SetPosition(servoStartPos)
	if err := servo.SendData(); err != nil {
		log.Println("unable to send servo data:", err)
	}
	log.Println("waiting for the servo")
	time.Sleep(time.Second * 1) // to be sure that the servo is on the right position
	lidar := Lidar{
		RPM:  lidarRPM,
		Mode: lidarMode,
		Process: Process{
			Args: []string{lidarPort, "--rpm", fmt.Sprint(lidarRPM), "--mode", fmt.Sprint(lidarMode)},
			Path: lidarExe, // TODO: Check if exists
		},
	}

	// Create communication channels
	lidarChan := make(chan *LidarCloud) // LidarCloud is >64kB so it cannot be directly passed by a channel
	servoChan := make(chan ServoData)
	accelChan := make(chan AccelDataUnion)

	// Create data buffers
	var lidarBuffer *LidarCloud
	accelBuffer := NewAccelDataBuffer(32)
	servoBuffer := NewServoDataBuffer(32) // unused

	// Goroutines
	go accel.StartLoop(accelChan)
	lidarStarted := false
	servoStarted := false

	// Fusion
	fusion := Fusion{
		CloudRotation: cloudRotation,
	}

	// Main loop
	for {
		select {
		case lidarData := <-lidarChan:
			lidarBuffer = lidarData
			fusion.Update(lidarBuffer, &accelBuffer)
		case servoData := <-servoChan:
			servoBuffer.Append(servoData)
		case accelData := <-accelChan:
			// when the accel is ready and its first measurement is read, start the lidar
			if !lidarStarted {
				go lidar.StartLoop(lidarChan)
				lidarStarted = true
			}
			if !servoStarted {
				go servo.StartLoop(servoChan)
				servoStarted = true
			}
			accelBuffer.Append(accelData)
		}
		writer.Flush()
	}
}
