package main

import (
	"errors"
	"log"
	"time"

	"github.com/knei-knurow/lidar-tools/frames"
)

type AccelData struct {
	xAccel int
	yAccel int
	zAccel int
	xGyro  int
	yGyro  int
	zGyro  int
	timept time.Time
}

var (
	accelCalib = AccelData{
		xAccel: -844,
		yAccel: 78,
		zAccel: 1542,
		xGyro:  244,
		yGyro:  -228,
		zGyro:  161,
	}
)

func calibrate(data *AccelData, calib *AccelData) {
	data.xAccel += calib.xAccel
	data.yAccel += calib.yAccel
	data.zAccel += calib.zAccel
	data.xGyro += calib.xGyro
	data.yGyro += calib.yGyro
	data.zGyro += calib.zGyro
}

func processAccelFrame(frame *frames.Frame) (AccelData, error) {
	timept := time.Now()
	var data AccelData

	if frame.Header()[0] != 'L' || frame.Header()[1] != 'D' || frame.Header()[2] != '-' {
		return data, errors.New("bad frame header")
	}

	if crc := frames.CalculateChecksum(*frame); crc != frame.Checksum() {
		// yeah but there is no checksum – for now, just print it
		// return data, errors.New("bad checksum")
		log.Println("bad checksum")
	}

	// TODO: make sure the lines below work correctly
	fdata := frame.Data()
	data.xAccel = mergeBytes(fdata[0], fdata[1])
	data.yAccel = mergeBytes(fdata[2], fdata[3])
	data.zAccel = mergeBytes(fdata[4], fdata[5])
	data.xGyro = mergeBytes(fdata[6], fdata[7])
	data.yGyro = mergeBytes(fdata[8], fdata[9])
	data.zGyro = mergeBytes(fdata[10], fdata[11])
	data.timept = timept

	calibrate(&data, &accelCalib)
	return data, nil
}

func mergeBytes(left8 byte, right8 byte) int {
	v := int((uint16(left8) << 8) | uint16(right8))
	// awesome conversion to signed int
	if v >= 32768 {
		v = -65536 + v
	}
	return v
}
