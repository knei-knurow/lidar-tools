package main

import (
	"errors"
	"io"
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
		// POSSIBLE ERROR SOURCE: Values differ depending on the temperature
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

func readAceelFrame(port io.Reader, data []byte, header rune) (err error) {
	scan := false
	for i := 0; i < 18; i++ {
		buf := make([]byte, 1)
		_, err := port.Read(buf)
		if err != nil {
			return errors.New("cannot read from port")
		}

		if scan {
			data[i] = buf[0]
		} else {
			if buf[0] == byte(header) {
				scan = true
				data[i] = buf[0]
			} else {
				return errors.New("lost some data")
			}
		}
	}
	return nil
}

func processAccelFrame(frame frames.Frame) (AccelData, error) {
	timept := time.Now()
	var data AccelData

	if frame[0] != frames.LidarHeader[0] ||
		frame[1] != frames.LidarHeader[1] ||
		frame[2] != 12 ||
		frame[3] != '+' {
		return data, errors.New("bad frame begin")
	}

	if !frames.Verify(frame) {
		return data, errors.New("bad checksum")
	}

	// TODO: make sure the lines below work correctly
	fdata := frame.Data()
	data.timept = timept // POSSIBLE ERROR SOURCE: Time of data receipt
	data.xAccel = mergeBytes(fdata[0], fdata[1])
	data.yAccel = mergeBytes(fdata[2], fdata[3])
	data.zAccel = mergeBytes(fdata[4], fdata[5])
	data.xGyro = mergeBytes(fdata[6], fdata[7])
	data.yGyro = mergeBytes(fdata[8], fdata[9])
	data.zGyro = mergeBytes(fdata[10], fdata[11])

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
