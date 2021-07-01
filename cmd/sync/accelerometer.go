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
	xAccel float32
	yAccel float32
	zAccel float32
	xGyro  float32
	yGyro  float32
	zGyro  float32
	timept time.Time
}

// AccelDataDMP contains accel data processed by Digital Motion Processor (quaternions)
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

// Supported data processing modes
const (
	AccelModeRaw = iota
	AccelModeDMP
)

// MPU-6050 constants. More details in the product documentation.
const (
	AccelScale2   = 16384.0
	AccelScale4   = 8192.0
	AccelScale8   = 4096.0
	AccelScale16  = 2048.0
	GyroScale250  = 131.0
	GyroScale500  = 65.5
	GyroScale1000 = 32.8
	GyroScale2000 = 16.4
)

// lidar-avr settings
const (
	AccelScaleDefault = AccelScale2
	GyroScaleDefault  = GyroScale250
	DeltaTimeDefault  = 0.02 // time in seconds between two measurements
)

type Accel struct {
	mode        int
	calibration AccelData
	accelScale  float32
	gyroScale   float32
	deltaTime   float32
	frequency   float32
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
		accel.PreprocessData()
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
	accel.data.raw.xAccel = float32(mergeBytes(fdata[0], fdata[1]))
	accel.data.raw.yAccel = float32(mergeBytes(fdata[2], fdata[3]))
	accel.data.raw.zAccel = float32(mergeBytes(fdata[4], fdata[5]))
	accel.data.raw.xGyro = float32(mergeBytes(fdata[6], fdata[7]))
	accel.data.raw.yGyro = float32(mergeBytes(fdata[8], fdata[9]))
	accel.data.raw.zGyro = float32(mergeBytes(fdata[10], fdata[11]))

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

// PreprocessData converts raw accel data to X * gravitational_acceleration and gyro to deg/s
func (accel *Accel) PreprocessData() {
	accel.data.raw.xAccel = (accel.data.raw.xAccel + accel.calibration.xAccel) / accel.accelScale
	accel.data.raw.yAccel = (accel.data.raw.yAccel + accel.calibration.yAccel) / accel.accelScale
	accel.data.raw.zAccel = (accel.data.raw.zAccel + accel.calibration.zAccel) / accel.accelScale
	accel.data.raw.xGyro = (accel.data.raw.xGyro + accel.calibration.xGyro) / accel.gyroScale
	accel.data.raw.yGyro = (accel.data.raw.yGyro + accel.calibration.yGyro) / accel.gyroScale
	accel.data.raw.zGyro = (accel.data.raw.zGyro + accel.calibration.zGyro) / accel.gyroScale
}

// mergeBytes merges two bytest to int
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
