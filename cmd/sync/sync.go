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
	avrPortName   string
	avrBaudRate   int
	lidarPortName string
	lidarMode     int
	lidarRPM      int
	lidarExe      string
	accelOut      bool
	servoOut      bool
	lidarOut      bool
	syncMode      int
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

	flag.BoolVar(&accelOut, "accel", false, "print accelerometer data on stdout")
	flag.BoolVar(&servoOut, "servo", false, "print set servo position on stdout")
	flag.BoolVar(&lidarOut, "lidar", false, "print lidar data on stdout")

	flag.IntVar(&syncMode, "mode", 0, "synchronization mode")

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
		data:       ServoData{positon: servoStartPos},
		positonMin: servoMinPos,
		positonMax: servoMaxPos,
		vector:     30,
		port:       port,
		delayMs:    60,
	}
	log.Println("setting the servo to the start position")
	servo.SetPosition(servoStartPos)
	log.Println("waiting for the servo")
	time.Sleep(time.Second) // to be sure that the servo is on the right position

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

	// Create data buffers
	accelBuffer := NewAccelDataBuffer(32)
	servoBuffer := NewServoDataBuffer(32)
	var lidarBuffer *LidarCloud

	if syncMode == 0 {
		for {
			servo.Move()
			if err := servo.SendData(); err != nil {
				log.Println("Unable to send servo data: ", err)
			}
			time.Sleep(time.Millisecond * 1000)
			// serv := float64(servo.data.positon-servoMinPos) * 0.05
			// <-lidarChan
			// lidarData := <-lidarChan
			// for i := 0; i < int(lidarData.Size); i++ {
			// 	pt := lidarData.Data[i]
			// 	x := float64(pt.Dist) * math.Cos(serv*math.Pi/180) * math.Sin(float64(pt.Angle)*math.Pi/180)
			// 	y := float64(pt.Dist) * math.Sin(serv*math.Pi/180) * math.Sin(float64(pt.Angle)*math.Pi/180)
			// 	z := float64(pt.Dist) * math.Cos(float64(pt.Angle)*math.Pi/180)
			// 	if pt.Dist != 0 {
			// 		writer.WriteString(fmt.Sprintf("L %f %f %f\n", x, y, z))
			// 	}
			// }
			// time.Sleep(time.Millisecond * 200)
		}

	}

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
			lidarBuffer = lidarData
		case servoData := <-servoChan:
			if servoOut {
				writer.WriteString(fmt.Sprintf("S %d %d\n", servoData.timept.UnixNano(), servoData.positon))
			}
			servoBuffer.Append(servoData)
		case accelData := <-accelChan:
			if accelOut {
				writer.WriteString(fmt.Sprintf("A %d\t%d\t%d\t%d\t%d\t%d\t%d\n", accelData.timept.UnixNano(),
					accelData.xAccel, accelData.yAccel, accelData.zAccel,
					accelData.xGyro, accelData.yGyro, accelData.zGyro))
			}
			accelBuffer.Append(accelData)
		}
		writer.Flush()

		go mergerLidarServoV1(lidarBuffer, &servoBuffer, true)
	}
}
