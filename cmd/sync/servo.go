package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/knei-knurow/lidar-tools/frames"
)

type Servo struct {
	positon    uint16    // given, not real
	positonMax uint16    // max position
	positonMin uint16    // min position
	vector     uint16    //
	timept     time.Time // time of sending a new position to the servo
	port       io.Writer // port to write controlling frames
}

// Move sends the move order to the servo and updates its movement vector.
func (servo *Servo) Move() {
	servo.positon += servo.vector

	switch {
	case servo.positon < servo.positonMin:
		servo.positon = servo.positonMin
		servo.vector = -servo.vector
	case servo.positon > servo.positonMax:
		servo.positon = servo.positonMax
		servo.vector = -servo.vector
	}
}

func (servo *Servo) SendData() (err error) {
	inputByte := servo.positon
	data := []byte{byte(inputByte >> 8), byte(inputByte)} // TODO: Check whether correct
	f := frames.Create([]byte(frames.LidarHeader), data)
	for i, currentByte := range f {
		if _, err := servo.port.Write([]byte{currentByte}); err != nil {
			return fmt.Errorf("cannot write data (byte %d) to port: %s", i, err)
		}
	}

	// POSSIBLE ERROR SOURCE: that's the time of the frame sending,
	// not actual servo set time
	servo.timept = time.Now()
	return nil
}

// StartLoop starts a loop responsible for controlling the servo position
// and updading the channel with its new calculated position.
func (servo *Servo) StartLoop(delay uint) {
	for {
		servo.Move()
		if err := servo.SendData(); err != nil {
			log.Println("Unable to send servo data: ", err)
		}

		if delay != 0 {
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
	}
}
