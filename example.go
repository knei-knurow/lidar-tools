package main

import (
	"fmt"

	"github.com/knei-knurow/lidar-tools/frames"
)

func main() {
	f := [][]byte{
		{'k', 'n', 'e', 'i'},
		{'d', 'o', 'n', 'd', 'u'},
		{'d', 'u', 'p', 'c', 'i', 'a'},
	}

	for i, data := range f {
		// Demostration of frames.CreateFrame function.
		frame := frames.Create([]byte(frames.LidarHeader), data)
		fmt.Printf("--- %d\n", i)
		fmt.Printf("data: % x\n", string(data))
		fmt.Printf("frame: % x\n", string(frame))
		fmt.Printf("funcs: header: % x, data: % x, checksum: %02x\n", frame.Header(), frame.Data(), frame.Checksum())
		for j, b := range frame {
			// Demonstration of frames.DescribeByte function.
			// It's very useful when you want to quickly examine what a particular byte represents.
			fmt.Printf("%2d: %s\n", j, frames.DescribeByte(b))
		}
	}
}
