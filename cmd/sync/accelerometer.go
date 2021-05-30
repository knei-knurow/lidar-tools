package main

import (
	"errors"
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

	if frame.Header[0] != 'L' || frame.Header[1] != 'D' || frame.Header[2] != '-' {
		return data, errors.New("bad frame header")
	}

	if crc := frames.CalculateCRC(frame.Data); crc != frame.Checksum {
		// yeah but there is no checksum
		// return data, errors.New("bad checksum")
	}

	data.xAccel = mergeBytes(frame.Data[0], frame.Data[1])
	data.yAccel = mergeBytes(frame.Data[2], frame.Data[3])
	data.zAccel = mergeBytes(frame.Data[4], frame.Data[5])
	data.xGyro = mergeBytes(frame.Data[6], frame.Data[7])
	data.yGyro = mergeBytes(frame.Data[8], frame.Data[9])
	data.zGyro = mergeBytes(frame.Data[10], frame.Data[11])
	data.timept = timept

	calibrate(&data, &accelCalib)
	return data, nil
}

func mergeBytes(left8 byte, right8 byte) (v int) {
	v = int((uint16(left8) << 8) | uint16(right8))
	// awesome conversion to signed int
	if v >= 32768 {
		v = -65536 + v
	}
	return
}
