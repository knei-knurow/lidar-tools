package main

import (
	"strconv"
	"strings"
)

func getCloudData(line string) (cloudIndex int, elapsed int, err error) {
	line = strings.TrimSuffix(line, "\n")

	splits := strings.Split(line, " ")

	cloudIndexStr := splits[1]
	cloudIndex, err = strconv.Atoi(cloudIndexStr)
	if err != nil {
		return
	}

	elapsedStr := splits[2]
	elapsed, err = strconv.Atoi(elapsedStr)
	if err != nil {
		return
	}

	return
}
