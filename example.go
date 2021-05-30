package main

import (
	"fmt"

	"github.com/knei-knurow/lidar-tools/frames"
)

func main() {
	dst := make([]byte, 5, 10)
	src := []byte{'d', 'u', 'p', 'c', 'i', 'a'}

	i := copy(dst, src)

	fmt.Println("copied:", i)
	fmt.Printf("dst: %s\n", dst)
	fmt.Printf("src: %s\n", src)

	fmt.Printf("\n- - -\n\n")

	// Demonstration of frames.DescribeByte function.
	// It's very useful when you want to quickly examine what a particular byte represents.
	// See its implementation to learn few cool tricks about formatting verbs.
	/*
		for i := 0; i < 8; i++ {
			fmt.Println(frames.DescribeByte(byte(i)))
		}
	*/

	fmt.Printf("\n- - -\n\n")

	// f1 and f2 are the same
	f1 := make(frames.Frame, 10)
	f2 := make(frames.Frame, 10)
	fmt.Println("f1:", f1, "f2:", f2)

	fmt.Printf("\n- - -\n\n")

	// var i uint16
	// for i = 0; i < 5; i++ {
	// 	rf := frames.EncodeRawFrame(i)
	// 	f := frames.EncodeFrame(i)
	// 	crc := frames.CalculateCRC(rf)
	// 	printFrame(rf, f, crc, int(i))
	// }
}

func printFrame(rf []byte, f []byte, crc byte, i int) {
	fmt.Println("---")
	fmt.Printf("raw frame %d: % x\n", i, rf)
	fmt.Printf("    frame %d: % x\n", i, f)
	fmt.Printf("    crc   %d: %s\n", i, frames.DescribeByte(crc))
}
