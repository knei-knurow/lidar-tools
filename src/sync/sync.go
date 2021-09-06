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

	// Servo flags
	servoStep  uint
	servoDelay uint
	servoMin   uint
	servoCalib uint
	servoStart uint
	servoMax   uint

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

	// Servo args
	flag.UintVar(&servoStep, "servostep", 2, "single servo step size")
	flag.UintVar(&servoDelay, "servodelay", 40, "delay in ms between steps")
	flag.UintVar(&servoMin, "servomin", servoMinPos, "min servo pos (might be corrected by AVR software)")
	flag.UintVar(&servoCalib, "servocalib", servoCalibPos, "servo position for accel calib (most horizontal position)")
	flag.UintVar(&servoStart, "servostart", servoMaxPos, "servo position for scan start")
	flag.UintVar(&servoMax, "servomax", servoMaxPos, "max servo pos (might be corrected by AVR software)")

	// Misc args
	flag.Float64Var(&cloudRotation, "cloudrotation", PrototypeCloudRotation, "each scanned 2D cloud will be rotated by CloudRotation radians, this value depends on the physical lidar location")

	flag.Parse()
	log.Println("starting...")
}

func main() {
	log.Println(cloudRotation)

	writer := bufio.NewWriter(os.Stdout)

	log.Println("opening AVR port")
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
		calibration: noAccelCalib,
		accelScale:  AccelScaleDefault,
		gyroScale:   GyroScaleDefault,
		deltaTime:   DeltaTimeDefault,
		port:        port,
		mode:        AccelModeRaw,
	}
	servo := Servo{
		data:       ServoData{positon: uint16(servoCalib)},
		positonMin: uint16(servoMin),
		positonMax: uint16(servoMax),
		vector:     uint16(servoStep),
		port:       port,
		delayMs:    servoDelay,
	}
	log.Println("servo is setting to the calibration position")
	servo.SetPosition(servoCalibPos)
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
			// when the accel is ready and its first measurement is read, start the servo and lidar
			if !servoStarted {
				log.Println("servo is setting to the start position")
				servo.SetPosition(uint16(servoStart))
				if err := servo.SendData(); err != nil {
					log.Println("unable to send servo data:", err)
				}
				log.Println("waiting for the servo")
				time.Sleep(time.Second * 2) // to be sure that the servo is on the right position

				go servo.StartLoop(servoChan)
				servoStarted = true
			}
			if !lidarStarted {
				go lidar.StartLoop(lidarChan)
				lidarStarted = true
			}
			accelBuffer.Append(accelData)
		}
		writer.Flush()
	}
}
