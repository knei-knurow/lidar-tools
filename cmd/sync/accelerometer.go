package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"time"

	"github.com/knei-knurow/frames"
)

// AccelData contains raw accel data
type AccelData struct {
	xAccel int
	yAccel int
	zAccel int
	xGyro  int
	yGyro  int
	zGyro  int
	timept time.Time
}

// AccelDataDMP contains accel data processed by Digital Motion Processor
type AccelDataDMP struct {
	qw     float32
	qx     float32
	qy     float32
	qz     float32
	timept time.Time
}

// AccelDataUnion is union-like structure which is used for sending accel data over channels
type AccelDataUnion struct {
	raw AccelData
	dmp AccelDataDMP
}

const (
	AccelModeRaw = iota
	AccelModeDMP
)

type Accel struct {
	mode        int
	calibration AccelData
	port        io.Reader
	data        AccelDataUnion
}

var (
	accelCalib = AccelData{
		// POSSIBLE ERROR SOURCE: Values differ depending on the temperature
		xAccel: 0,
		yAccel: 0,
		zAccel: 0,
		xGyro:  0,
		yGyro:  0,
		zGyro:  0,
	}
)

// StartLoop starts the accelerometer main loop
func (accel *Accel) StartLoop(channel chan AccelDataUnion) {
	for {
		accel.ReadData()
		channel <- accel.data
	}
}

func (accel *Accel) ProcessAccelFrame(frame frames.Frame) (err error) {
	timept := time.Now()

	if frame[0] != 'L' ||
		frame[1] != 'D' ||
		frame[2] != 12 ||
		frame[3] != '+' {
		return errors.New("bad frame begin")
	}

	if !frames.Verify(frame) {
		return errors.New("bad checksum")
	}

	fdata := frame.Data()
	accel.data.raw.timept = timept // POSSIBLE ERROR SOURCE: Time of data receipt
	accel.data.raw.xAccel = mergeBytes(fdata[0], fdata[1])
	accel.data.raw.yAccel = mergeBytes(fdata[2], fdata[3])
	accel.data.raw.zAccel = mergeBytes(fdata[4], fdata[5])
	accel.data.raw.xGyro = mergeBytes(fdata[6], fdata[7])
	accel.data.raw.yGyro = mergeBytes(fdata[8], fdata[9])
	accel.data.raw.zGyro = mergeBytes(fdata[10], fdata[11])

	accel.calibrate()
	return nil
}

func (accel *Accel) ProcessAccelFrameDMP(frame frames.Frame) (err error) {
	timept := time.Now()

	if frame[0] != 'L' ||
		frame[1] != 'Q' ||
		frame[2] != 16 ||
		frame[3] != '+' {
		return errors.New("bad frame begin")
	}

	if !frames.Verify(frame) {
		return errors.New("bad checksum")
	}

	fdata := frame.Data()
	accel.data.dmp.timept = timept // POSSIBLE ERROR SOURCE: Time of data receipt
	accel.data.dmp.qw = float32frombytes(fdata[0:4])
	accel.data.dmp.qx = float32frombytes(fdata[4:8])
	accel.data.dmp.qy = float32frombytes(fdata[8:12])
	accel.data.dmp.qz = float32frombytes(fdata[12:16])

	return nil
}

func (accel *Accel) calibrate() {
	accel.data.raw.xAccel += accel.calibration.xAccel
	accel.data.raw.yAccel += accel.calibration.yAccel
	accel.data.raw.zAccel += accel.calibration.zAccel
	accel.data.raw.xGyro += accel.calibration.xGyro
	accel.data.raw.yGyro += accel.calibration.yGyro
	accel.data.raw.zGyro += accel.calibration.zGyro
}

// mergeBytes Merges two bytest to int
func mergeBytes(left8 byte, right8 byte) int {
	v := int((uint16(left8) << 8) | uint16(right8))
	// awesome conversion to signed int
	if v >= 32768 {
		v = -65536 + v
	}
	return v
}

// float32frombytes converts 4 bytes to float32
func float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

// ReadData reads and parses new measurement
func (accel *Accel) ReadData() (err error) {
	var frame frames.Frame
	var frameLen int
	if accel.mode == AccelModeRaw {
		frameLen = 18
	} else {
		frameLen = 22
	}

	frame = make(frames.Frame, frameLen)
	if err := accel.ReadAccelFrame(frame, frameLen); err != nil {
		log.Printf("error: %s\n", err)
	}

	if accel.mode == AccelModeRaw {
		err = accel.ProcessAccelFrame(frame)
	} else {
		err = accel.ProcessAccelFrameDMP(frame)
	}
	if err != nil {
		return errors.New("cannot process frame")
	}
	return nil
}

// ReadAccelFrame is a low level function to read an accelerometer frame
func (accel *Accel) ReadAccelFrame(data []byte, length int) (err error) {
	scan := false
	for i := 0; i < length; i++ {
		buf := make([]byte, 1)
		_, err := accel.port.Read(buf)
		if err != nil {
			return errors.New("cannot read from port")
		}

		if scan {
			data[i] = buf[0]
		} else {
			if buf[0] == 'L' {
				scan = true
				data[i] = buf[0]
			} else {
				return errors.New("lost some data")
			}
		}
	}
	return nil
}
