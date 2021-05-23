// Package frame provides standard functions to deal with
// frames used in lidar-related tools.
package frame

import (
	"fmt"
	"strings"
)

const (
	FrameLidar  = "LD"
	FrameMotors = "MT"
)

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
		// fmt.Printf("crc %d: % x\n", i, crc)
	}
	return
}

/// DescribeByte prints everything that a single byte can represent.
/// Ir prints binary value, decimal value and ASCII character.
func DescribeByte(b byte) string {
	return fmt.Sprintf("byte (bin: %b, dec: %d, ASCII: %q)", b, b, b)
}
