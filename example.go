package main

import (
	"fmt"

	"github.com/knei-knurow/frames"
)

var testCases = []struct {
	inputHeader      [2]byte
	inputData        []byte
	expectedChecksum byte
	expectedLength   int
}{
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{},
		expectedChecksum: 0x00,
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'A'},
		expectedChecksum: 0x41,
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'t', 'e', 's', 't'},
		expectedChecksum: 0x16,
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'d', 'u', 'p', 'c', 'i', 'a'},
		expectedChecksum: 0x0a,
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'l', 'o', 'l', 'x', 'd'},
		expectedChecksum: 0x73,
	},
	{
		inputHeader:      [2]byte{'M', 'T'},
		inputData:        []byte{'d', 'o', 'n', 'd', 'u'},
		expectedChecksum: 0x30,
	},
	// Invalid: header must always be 2 bytes length, data equal or more than 1 byte
	// {
	// 	inputHeader:      []byte{},
	// 	inputData:        []byte{},
	// 	expectedChecksum: 0x08,
	// },
}

func main() {
	// f := [][]byte{
	// 	{'k', 'n', 'e', 'i'},
	// 	{'d', 'o', 'n', 'd', 'u'},
	// 	{'d', 'u', 'p', 'c', 'i', 'a'},
	// }

	for i, tc := range testCases {
		// Demostration of frames.CreateFrame function.
		frame := frames.Create(tc.inputHeader, tc.inputData)
		fmt.Printf("--- %d\n", i)
		fmt.Printf("data: % x (%q)\n", string(tc.inputData), string(tc.inputData))
		fmt.Printf("frame: % x\n", string(frame))
		fmt.Printf("funcs: header: % x, data: % x, checksum: %02x\n", frame.Header(), frame.Data(), frame.Checksum())
		for j, b := range frame {
			// Demonstration of frames.DescribeByte function.
			// It's very useful when you want to quickly examine what a particular byte represents.
			fmt.Printf("%2d: %s\n", j, frames.DescribeByte(b))
		}
	}
}
