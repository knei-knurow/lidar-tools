// Package frames provides standard functions to deal with
// frames used in lidar-related tools.
package frames

import (
	"fmt"
)

// FrameHeader represents the type of a frame. It takes the form of 2 uppercase ASCII characters.
type FrameHeader string

const (
	FrameLidar  FrameHeader = "LD" // Frame format used for lidar-related stuff.
	FrameMotors FrameHeader = "MT" // Frame format used for motors-related stuff.
)

// Frame represents a frame that can be e.g sent by USART.
//
// Header is always 2 bytes, CRC is always 1 byte.
//
// Example frame (H = header byte, D = data byte, C = CRC byte):
//
// BB+DDDDD#C
type Frame []byte

// Header returns frame's 2 leading bytes.
func (f Frame) Header() []byte {
	return f[0:2]
}

// SetHeader sets frame's header.
func (f Frame) SetHeader(header [2]byte) {
	f[0] = header[0]
	f[1] = header[1]
}

// SetData sets frame's data and recalculates CRC.
// TODO: discuss how this should work.
// func (f Frame) SetData(data []byte) {
// 	// allocate enough space for data
// 	if len(f)-3 > len(data) {
// 		f = copy()
// 	}

// 	f = "MT+"

// 	for i, b := range data {
// 		placeAt := 3 + i
// 		if placeAt > len(f) {
// 			f = append(f, )
// 			copy
// 		}
// 		f[] = b
// 	}

// 	// recalculate CRC
// }

// Header returns frame's part from fourth to antepenultimate byte.
func (f Frame) Data() []byte {
	return f[3 : len(f)-3]
}

// Header returns frame's last byte - a simple CRC checksum.
func (f Frame) Checksum() byte {
	return f[len(f)-1]
}

// CalculateChecksum calculates the CRC checksum of frame.
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
