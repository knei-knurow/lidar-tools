package main

import (
	"fmt"

	"github.com/knei-knurow/lidar-tools/frames"
)

func main() {
	// Demonstration of frame.CalculateCRC and frame.DescribeByte functions.
	data := []byte{123, 153, 223}
	fmt.Printf("data: %03b, crc: %s\n", data, frames.DescribeByte(frames.CalculateCRC(data)))

	// Demonstration of frames.DescribeByte function.
	// It's very useful to quickly examine what a particular byte means.
	// See its implementation to learn few cool tricks about formatting verbs.

	fmt.Printf("\n- - -\n\n")
	for i := 0; i < 8; i++ {
		fmt.Println(frames.DescribeByte(byte(i)))
	}

	fmt.Printf("\n- - -\n\n")

	// f1 and f2 are the same
	f1 := frames.Frame{Header: frames.FrameLidar, Data: []byte{48, 48}, Checksum: 0}
	f2 := frames.Frame{Header: frames.FrameLidar, Data: []byte("00"), Checksum: 0}
	fmt.Println("f1:", f1, "f2:", f2)

	fmt.Printf("\n- - -\n\n")

	var i uint16
	for i = 0; i < 5; i++ {
		rf := frames.EncodeRawFrame(i)
		f := frames.EncodeFrame(i)
		crc := frames.CalculateCRC(rf)
		printFrame(rf, f, crc, int(i))
	}
}

func printFrame(rf []byte, f []byte, crc byte, i int) {
	fmt.Println("---")
	fmt.Printf("raw frame %d: % x\n", i, rf)
	fmt.Printf("    frame %d: % x\n", i, f)
	fmt.Printf("    crc   %d: %s\n", i, frames.DescribeByte(crc))
}
