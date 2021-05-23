// Package frame provides standard functions to deal with
// frames used in lidar-related tools.
package frame

import (
	"fmt"
	"strings"
)

type FrameHeader string

const (
	FrameLidar  FrameHeader = "LD"
	FrameMotors FrameHeader = "MT"
)

// Frame is a standard frame used in the rover project.
type Frame struct {
	Header   FrameHeader
	Data     []byte
	Checksum byte
}

func (f Frame) String() string {
	return fmt.Sprintf("%s+%x#%x", f.Header, f.Data, f.Checksum)
}

// EncodeRawFrame creates a frame transporting data.
// It does not have CRC checksum.
func EncodeRawFrame(data uint16) (frame []byte) {
	var builder strings.Builder
	builder.Grow(6)

	builder.WriteString("LD+")
	builder.WriteByte(byte(data >> 8)) // Write most significant 8 bits
	builder.WriteByte(byte(data))      // Write least significant 8 bits
	builder.WriteString("#")

	frame = []byte(builder.String())
	return
}

// EncodeFrame creates a standard frame transporting data.
func EncodeFrame(data uint16) (frame []byte) {
	var builder strings.Builder
	builder.Grow(2)

	rawFrame := EncodeRawFrame(data)
	builder.WriteString(string(rawFrame))

	crc := CalculateCRC([]byte(rawFrame))
	builder.WriteByte(crc)

	frame = []byte(builder.String())
	return
}

// CalculateCRC calculates the CRC checksum of data.
func CalculateCRC(data []byte) (crc byte) {
	crc = data[0]
	for i := 1; i < len(data); i++ {
		crc ^= data[i]
	}
	return
}

// DescribeByte prints everything most common representations of a byte.
// It prints b's binary value, decimal, hexadecimal value and ASCII.
func DescribeByte(b byte) string {
	return fmt.Sprintf("byte(bin: %08b, dec: %3d, hex: %02x, ASCII: %+q)", b, b, b, b)
}
