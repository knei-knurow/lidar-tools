package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	rps = 660 / 60
	pts = 1000
)

func main() {
	stdout := bufio.NewWriter(os.Stdout)

	cnt := 0
	timeStart := time.Now()
	stdout.WriteString("! 0 0\n")
	for {
		cnt++
		timeDiff := time.Now().Sub(timeStart)

		stdout.WriteString("# Dummy lidar-scan data\n") // a comment line for tests
		for i := 0; i < 1000; i++ {
			stdout.WriteString(fmt.Sprintf("%f %f\n", float32(i)/float32(pts)*360, float32(cnt*10000)+float32(i)))
		}
		stdout.WriteString(fmt.Sprintf("! %d %d\n", cnt, timeDiff.Milliseconds()))
		stdout.Flush()

		timeStart = time.Now()

		time.Sleep(time.Millisecond * (1000 / rps))
	}
}
