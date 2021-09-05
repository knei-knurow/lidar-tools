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
	avrPortName string
	avrBaudRate int

	lidarPortName string
	lidarMode     int
	lidarRPM      int
	lidarExe      string

	accelExe string

	accelOut bool
	estOut   bool
	servoOut bool
	lidarOut bool

	accelTest bool
	servoTest bool
	lidarTest bool
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("sync: ")

	flag.StringVar(&avrPortName, "avrport", "COM13", "AVR serial communication port")
	flag.IntVar(&avrBaudRate, "avrbaud", 19200, "port baud rate (bps)")

	flag.StringVar(&lidarExe, "lidarexe", "lidar.exe", "lidar-scan executable")
	flag.StringVar(&lidarPortName, "lidarport", "COM4", "RPLIDAR serial communication port")
	flag.IntVar(&lidarMode, "lidarmode", rplidarModeDefault, "RPLIDAR mode")
	flag.IntVar(&lidarRPM, "lidarpm", 660, "RPLIDAR given revolutions per minute")

	flag.StringVar(&accelExe, "accelexe", "attitude-estimator.exe", "attitude estimator executable")

	flag.BoolVar(&accelOut, "accel", false, "print accelerometer data on stdout")
	flag.BoolVar(&estOut, "est", false, "print attitude estimator data on stdout")
	flag.BoolVar(&servoOut, "servo", false, "print set servo position on stdout")
	flag.BoolVar(&lidarOut, "lidar", false, "print lidar data on stdout")

	flag.BoolVar(&accelTest, "acceltest", false, "perform accelerometer test and exit (to check connection, power, etc)")
	flag.BoolVar(&servoTest, "servotest", false, "perform servo test and exit (to check connection, power, etc)")
	flag.BoolVar(&lidarTest, "lidartest", false, "perform lidar test and exit (to check connection, power, etc)")

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

	lidar := Lidar{ // TODO: make it more configurable from command line
		RPM:  lidarRPM,
		Mode: lidarMode,
		Process: Process{
			Args: []string{lidarPortName, "--rpm", fmt.Sprint(lidarRPM), "--mode", fmt.Sprint(lidarMode)},
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
	servoBuffer := NewServoDataBuffer(32)

	// Goroutines
	go accel.StartLoop(accelChan)
	lidarStarted := false
	servoStarted := false

	// Fusion stuff
	var fusion Fusion

	// Main loop
	for {
		select {
		case lidarData := <-lidarChan:
			if lidarOut {
				writer.WriteString(fmt.Sprintf("L %d %d\n", lidarData.ID, lidarData.Size))
			}
			lidarBuffer = lidarData
			fusion.Update(lidarBuffer, &accelBuffer)
		case servoData := <-servoChan:
			if servoOut {
				writer.WriteString(fmt.Sprintf("S %d %d\n", servoData.timept.UnixNano(), servoData.positon))
			}
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

			if accelOut {
				if accel.mode == AccelModeRaw {
					writer.WriteString(fmt.Sprintf("A %d\t%f\t%f\t%f\t%f\t%f\t%f\n",
						accelData.raw.timept.UnixNano(),
						accelData.raw.xAccel, accelData.raw.yAccel, accelData.raw.zAccel,
						accelData.raw.xGyro, accelData.raw.yGyro, accelData.raw.zGyro))
				} else {
					writer.WriteString(fmt.Sprintf("a %d\t%f\t%f\t%f\t%f\n",
						accelData.quat.timept.UnixNano(),
						accelData.quat.qw, accelData.quat.qx, accelData.quat.qy, accelData.quat.qz))
				}
			}
			if estOut {
				writer.WriteString(fmt.Sprintf("%f\t%f\t%f\t%f\n",
					accelData.quat.qw, accelData.quat.qx, accelData.quat.qy, accelData.quat.qz))
			}
			accelBuffer.Append(accelData)
		}
		writer.Flush()
	}
}
