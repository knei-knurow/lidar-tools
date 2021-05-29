package main

import (
	"time"
)

type Servo struct {
	positon    uint16    // given, not real
	positonMax uint16    // max position
	positonMin uint16    // min position
	vector     uint16    //
	timept     time.Time // time of sending a new position to the servo
}

func (servo *Servo) move() {
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
