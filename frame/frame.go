// Package frame provides standard functions to deal with
// frames used in lidar-related tools.
package frame

import (
	"fmt"
	"strings"
)

// CreateFrame creates a standard frame transporting data
// and returns it.
func CreateFrame(data uint16) (frame []byte) {
	var builder strings.Builder
	builder.Grow(7)

	builder.WriteString("LD+")
	builder.WriteByte(byte(data >> 8))
	builder.WriteByte(byte(data))
	builder.WriteString("#")

	encoded := builder.String()
	crc := CalculateCRC([]byte(encoded))

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

/// DescribeByte prints everything that a single byte can represent.
/// Ir prints binary value, decimal value and ASCII character.
func DescribeByte(b byte) string {
	return fmt.Sprintf("byte (bin: %b, dec: %d, ASCII: %q)", b, b, b)
}
