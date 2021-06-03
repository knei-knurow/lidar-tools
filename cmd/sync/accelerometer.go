package main

import (
	"errors"
	"io"
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

type Accel struct {
	calibration AccelData
	port        io.Reader
	data        AccelData
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

// StartLoop starts the accelerometer main loop
func (accel *Accel) StartLoop() {
	for {
		accel.ReadData()
	}
}

func (accel *Accel) ProcessAccelFrame(frame frames.Frame) (err error) {
	timept := time.Now()

	if frame[0] != frames.LidarHeader[0] ||
		frame[1] != frames.LidarHeader[1] ||
		frame[2] != 12 ||
		frame[3] != '+' {
		return errors.New("bad frame begin")
	}

	if !frames.Verify(frame) {
		return errors.New("bad checksum")
	}

	// TODO: make sure the lines below work correctly
	fdata := frame.Data()
	accel.data.timept = timept // POSSIBLE ERROR SOURCE: Time of data receipt
	accel.data.xAccel = mergeBytes(fdata[0], fdata[1])
	accel.data.yAccel = mergeBytes(fdata[2], fdata[3])
	accel.data.zAccel = mergeBytes(fdata[4], fdata[5])
	accel.data.xGyro = mergeBytes(fdata[6], fdata[7])
	accel.data.yGyro = mergeBytes(fdata[8], fdata[9])
	accel.data.zGyro = mergeBytes(fdata[10], fdata[11])

	accel.calibrate()
	return nil
}

func (accel *Accel) calibrate() {
	accel.data.xAccel += accel.calibration.xAccel
	accel.data.yAccel += accel.calibration.yAccel
	accel.data.zAccel += accel.calibration.zAccel
	accel.data.xGyro += accel.calibration.xGyro
	accel.data.yGyro += accel.calibration.yGyro
	accel.data.zGyro += accel.calibration.zGyro
}

func mergeBytes(left8 byte, right8 byte) int {
	v := int((uint16(left8) << 8) | uint16(right8))
	// awesome conversion to signed int
	if v >= 32768 {
		v = -65536 + v
	}
	return v
}

// ReadData reads and parses new measurement
func (accel *Accel) ReadData() (err error) {
	frame := make(frames.Frame, 18)
	if err := accel.ReadAceelFrame(frame); err != nil {
		log.Printf("error: %s\n", err)
	}
	err = accel.ProcessAccelFrame(frame)
	if err != nil {
		return errors.New("cannot process frame")
	}
	return nil
}

// ReadAceelFrame is a low level function to read an accelerometer frame
func (accel *Accel) ReadAceelFrame(data []byte) (err error) {
	scan := false
	for i := 0; i < 18; i++ {
		buf := make([]byte, 1)
		_, err := accel.port.Read(buf)
		if err != nil {
			return errors.New("cannot read from port")
		}

		if scan {
			data[i] = buf[0]
		} else {
			if buf[0] == frames.LidarHeader[0] {
				scan = true
				data[i] = buf[0]
			} else {
				return errors.New("lost some data")
			}
		}
	}
	return nil
}
