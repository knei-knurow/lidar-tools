// Package frames provides standard functions to deal with
// data frames in rover project.
package frames

import (
	"fmt"
	"strings"
)

// FrameHeader represents the type of a frame. It takes the form of 2 uppercase ASCII characters.
// type FrameHeader string

// Standard frame headers. Must be ASCII-only strings.
const (
	// Frame format used for lidar-related stuff.
	LidarHeader = "LD"

	// Frame format used for motors-related stuff.
	MotorsHeader = "MT"
)

// Frame represents a frame that can be e.g sent by USART.
//
// Frame starts with an arbitrary-length header (thought it is almost always 2 bytes).
// After a header comes a plus sign ("+").
// Then comes an arbitrary-length data.
// Data is terminated with a hash sign ("#").
// The last byte is a simple 8-bit CRC checksum.
//
// Example frame (H = header byte, D = data byte, C = CRC byte):
//
// HH+DDDDD#C
type Frame []byte

// Header returns frame's header. It is usually 2 bytes.
func (f Frame) Header() []byte {
	end := strings.IndexByte(string(f), '+')
	return f[0:end]
}

// Data returns frame's data part from the first byte after a plus sign ("+") up
// to the antepenultimate (last but one - 1) byte.
func (f Frame) Data() []byte {
	start := strings.IndexByte(string(f), '+')
	return f[start+1 : len(f)-2]
}

// Checksum returns frame's last byte - a simple CRC checksum.
func (f Frame) Checksum() byte {
	return f[len(f)-1]
}

// Create creates a new frame.
// The frame starts with header and contains data.
// It also calculates the checksum using frames.CalculateChecksum.
func Create(header []byte, data []byte) (frame Frame) {
	frame = make(Frame, len(header)+1+len(data)+2)

	copy(frame[:len(header)], header)
	frame[len(header)] = '+'
	copy(frame[len(header)+1:len(frame)-2], data)
	frame[len(frame)-2] = '#'
	frame[len(frame)-1] = CalculateChecksum(frame)

	return
}

// Assemble creates a frame from already available values.
func Assemble(header []byte, data []byte, checksum byte) (frame Frame) {
	frame = make(Frame, len(header)+1+len(data)+2)

	copy(frame[:len(header)], header)
	frame[len(header)] = '+'
	copy(frame[len(header)+1:len(frame)-2], data)
	frame[len(frame)-2] = '#'
	frame[len(frame)-1] = checksum

	return
}

// CalculateChecksum calculates the simple CRC checksum of frame.
//
// It takes all frame's bytes into account, except the last byte, because
// the last byte is the CRC itself.
func CalculateChecksum(f Frame) (crc byte) {
	crc = f[0]
	for i := 1; i < len(f)-1; i++ {
		crc ^= f[i]
	}

	return
}

func (f Frame) String() string {
	return fmt.Sprintf("%s+%x#%x", f.Header(), f.Data(), f.Checksum())
}

// DescribeByte prints everything most common representations of a byte.
// It prints b's binary value, decimal, hexadecimal value and ASCII.
func DescribeByte(b byte) string {
	return fmt.Sprintf("byte(bin: %08b, dec: %3d, hex: %02x, ASCII: %+q)", b, b, b, b)
}
