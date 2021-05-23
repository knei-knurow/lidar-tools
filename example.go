package main

import (
	"fmt"

	"github.com/bartekpacia/lidar-tools/frame"
)

func main() {
	data := []byte{123, 153, 223}
	fmt.Printf("data: %03b, crc: %s\n", data, frame.DescribeByte(frame.CalculateCRC(data)))
	// Demonstration of frame.DescribeByte function.
	// It's very useful to quickly examine what a particular byte means.
	// See its implementation to learn few cool tricks about formatting verbs.
	/*
		for i := 0; i < 256; i++ {
			fmt.Println(frame.DescribeByte(byte(i)))
		}
	*/

	// f1 and f2 are the same
	f1 := frame.Frame{Header: frame.FrameLidar, Data: []byte{48, 48}, Checksum: 0}
	f2 := frame.Frame{Header: frame.FrameLidar, Data: []byte("00"), Checksum: 0}
	fmt.Println(f1, f2)

	// rf0 := frame.EncodeRawFrame(0)
	// f0 := frame.EncodeFrame(0)
	// crc0 := frame.CalculateCRC(rf0)
	// printFrame(rf0, f0, crc0, 0)

	var i uint16
	for i = 0; i < 5; i++ {
		rf := frame.EncodeRawFrame(i)
		f := frame.EncodeFrame(i)
		crc := frame.CalculateCRC(rf)
		printFrame(rf, f, crc, int(i))
	}
}

func printFrame(rf []byte, f []byte, crc byte, i int) {
	fmt.Println("---")
	fmt.Printf("raw frame %d: % x\n", i, rf)
	fmt.Printf("    frame %d: % x\n", i, f)
	fmt.Printf("    crc   %d: %s\n", i, frame.DescribeByte(crc))
}
