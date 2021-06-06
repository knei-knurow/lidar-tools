package main

import "errors"

type AccelDataBuffer struct {
	size   int
	data   []AccelData
	pos    int
	isFull bool
}

func NewAccelDataBuffer(size int) (buffer AccelDataBuffer) {
	buffer.size = size
	buffer.data = make([]AccelData, size)
	return buffer
}

func (buffer *AccelDataBuffer) Append(element AccelData) (err error) {
	if buffer.size == 0 {
		return errors.New("buffer size equals 0")
	}

	buffer.data[buffer.pos] = element

	buffer.pos = (buffer.pos + 1) % buffer.size
	if buffer.pos == 0 {
		buffer.isFull = true
	}

	return nil
}

func (buffer *AccelDataBuffer) Get(posFromTop int) (element AccelData, err error) {
	if posFromTop > buffer.size {
		return AccelData{}, errors.New("too high position")
	}

	pos := buffer.pos - 1 - posFromTop
	if pos < 0 {
		if !buffer.isFull {
			return AccelData{}, errors.New("too high position because buffer is not full")
		}
		pos += buffer.size
	}

	return buffer.data[pos], nil
}

type ServoDataBuffer struct {
	size   int
	data   []ServoData
	pos    int
	isFull bool
}

func NewServoDataBuffer(size int) (buffer ServoDataBuffer) {
	buffer.size = size
	buffer.data = make([]ServoData, size)
	return buffer
}

func (buffer *ServoDataBuffer) Append(element ServoData) (err error) {
	if buffer.size == 0 {
		return errors.New("buffer size equals 0")
	}

	buffer.data[buffer.pos] = element

	buffer.pos = (buffer.pos + 1) % buffer.size
	if buffer.pos == 0 {
		buffer.isFull = true
	}

	return nil
}

func (buffer *ServoDataBuffer) Get(posFromTop int) (element ServoData, err error) {
	if posFromTop > buffer.size {
		return ServoData{}, errors.New("too high position")
	}

	pos := buffer.pos - 1 - posFromTop
	if pos < 0 {
		if !buffer.isFull {
			return ServoData{}, errors.New("too high position because buffer is not full")
		}
		pos += buffer.size
	}

	return buffer.data[pos], nil
}
