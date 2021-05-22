package main

import (
	"strings"
)

func createFrame(input uint16) []byte {
	var builder strings.Builder
	builder.Grow(7)

	builder.WriteString("LD+")
	builder.WriteByte(byte(input >> 8))
	builder.WriteByte(byte(input))
	builder.WriteString("#")

	encoded := builder.String()
	crc := calculateCRC([]byte(encoded))

	builder.WriteByte(crc)

	return []byte(builder.String())
}

func calculateCRC(value []byte) byte {
	crc := value[0]

	for i := 1; i < len(value); i++ {
		crc ^= value[i]
	}

	return crc
}
