package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/knei-knurow/lidar-tools/frames"
)

type ServoData struct {
	positon uint16    // given, not real
	timept  time.Time // time of sending a new position to the servo
}

type Servo struct {
	data       ServoData // servo data
	positonMax uint16    // max position
	positonMin uint16    // min position
	vector     uint16    //
	port       io.Writer // port to write controlling frames
	delayMs    uint      // ms delay between orders
}

// Move sends the move order to the servo and updates its movement vector.
func (servo *Servo) Move() {
	servo.data.positon += servo.vector

	switch {
	case servo.data.positon < servo.positonMin:
		servo.data.positon = servo.positonMin
		servo.vector = -servo.vector
	case servo.data.positon > servo.positonMax:
		servo.data.positon = servo.positonMax
		servo.vector = -servo.vector
	}
}

func (servo *Servo) SendData() (err error) {
	inputByte := servo.data.positon
	data := []byte{byte(inputByte >> 8), byte(inputByte)} // TODO: Check whether correct
	f := frames.Create([]byte(frames.LidarHeader), data)
	for i, currentByte := range f {
		if _, err := servo.port.Write([]byte{currentByte}); err != nil {
			return fmt.Errorf("cannot write data (byte %d) to port: %s", i, err)
		}
	}

	// POSSIBLE ERROR SOURCE: that's the time of the frame sending,
	// not actual servo set time
	servo.data.timept = time.Now()
	return nil
}

// StartLoop starts a loop responsible for controlling the servo position
// and updading the channel with its new calculated position.
func (servo *Servo) StartLoop(channel chan ServoData) {
	for {
		servo.Move()
		if err := servo.SendData(); err != nil {
			log.Println("Unable to send servo data: ", err)
		}
		channel <- servo.data

		if servo.delayMs != 0 {
			time.Sleep(time.Millisecond * time.Duration(servo.delayMs))
		}
	}
}
