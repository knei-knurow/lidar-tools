package main

import (
	"errors"
	"time"

	"github.com/knei-knurow/lidar-tools/frame"
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

func processAccelFrame(frm *frame.Frame) (AccelData, error) {
	timept := time.Now()
	var data AccelData

	if frm.Header[0] != 'L' || frm.Header[1] != 'D' || frm.Header[2] != '-' {
		return data, errors.New("Bad frame header.")
	}

	if crc := frame.CalculateCRC(frm.Data); crc != frm.Checksum {
		// yeah but there is no checksum
		// return data, errors.New("Bad checksum.")
	}

	data.xAccel = mergeBytes(frm.Data[0], frm.Data[1])
	data.yAccel = mergeBytes(frm.Data[2], frm.Data[3])
	data.zAccel = mergeBytes(frm.Data[4], frm.Data[5])
	data.xGyro = mergeBytes(frm.Data[6], frm.Data[7])
	data.yGyro = mergeBytes(frm.Data[8], frm.Data[9])
	data.zGyro = mergeBytes(frm.Data[10], frm.Data[11])
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
