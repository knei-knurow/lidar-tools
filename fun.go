package main

import (
	"fmt"

	"github.com/bartekpacia/lidar-tools/frame"
)

func main() {

	// temporarily - just for testing
	fmt.Printf("\x10 in ascii: %q\n", 0)
	rf0 := frame.EncodeRawFrame(0)
	f0 := frame.EncodeFrame(0)
	crc0 := frame.CalculateCRC(rf0)
	fmt.Printf("rawFrame0: % x\n", rf0)
	fmt.Printf("frame0: % x\n", f0)
	fmt.Printf("crc0: % x\n", crc0)

	fmt.Println("---")
	rf1 := frame.EncodeRawFrame(3)
	f1 := frame.EncodeFrame(3)
	crc1 := frame.CalculateCRC(rf1)
	fmt.Printf("rawFrame1: % x\n", rf1)
	fmt.Printf("frame1: % x\n", f1)
	fmt.Printf("crc1: % x\n", crc1)

	fmt.Println("---")
	rf2 := frame.EncodeRawFrame(5)
	f2 := frame.EncodeFrame(5)
	crc2 := frame.CalculateCRC(rf2)
	fmt.Printf("rawFrame2: % x\n", rf2)
	fmt.Printf("frame2: % x\n", f2)
	fmt.Printf("crc2: % x\n", crc2)
}
