package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Represents a single angle+dist measurement
type Point struct {
	angle float32 // angle in degrees
	dist  float32 // distance in millimeters
	// pointpt time.Time // estimated measurement time point
}

// Contains one full-360deg lidar point cloud
type LidarCloud struct {
	id        int       //
	timeBegin time.Time // time point of starting line read (line starting with '!')
	timeDiff  int       // number of milliseconds of current cloud measurement (received from lidar-scan)
	timeEnd   time.Time // timeBegin increased by timeDiff milliseconds
	data      []Point   // measurements data
}

// rplidar scanning modes
const ( // rplidar modes
	rplidarModeBoost       = 2
	rplidarModeSensitivity = 3
	rplidarModeStability   = 4
	rplidarModeDefault     = rplidarModeSensitivity
)

// General lidar parameters
type Lidar struct {
	timeInit time.Time     // Time of the first starting line read (line starting with '!')
	rpm      int           // declared rpm (but actual may differ)
	mode     int           // rplidar mode
	process  *exec.Cmd     // lidar-scan process
	argv     []string      // lidar-scan process argv
	path     string        // lidar-scan path
	running  bool          // whether lidar-scan is currently scanning
	stdout   io.ReadCloser // lidar-scan stdout
	stderr   io.ReadCloser // lidar-scan stderr
}

// Starts the lidar-scan process.
// Does not check whether it has been already started.
func (lidar *Lidar) processStart() error {
	var err error
	lidar.process = exec.Command(lidar.path, lidar.argv...)
	log.Println("starting lidar-scan process")

	lidar.stdout, err = lidar.process.StdoutPipe()
	if err != nil {
		return errors.New("Unable to get stdout of lidar-scan process.")
	}
	lidar.stderr, err = lidar.process.StderrPipe()
	if err != nil {
		return errors.New("Unable to get stderr of lidar-scan process.")
	}

	err = lidar.process.Start()
	if err != nil {
		return errors.New("Unable to start lidar-scan process.")
	}

	lidar.running = true
	return nil
}

// Sends interrupt signal (ctrl+c) to the lidar-scan process,
// so it should be able to handle it and perform cleanup.
// Does not check whether it has been already started.
// On Windows has the same bahaviour like processKill.
func (lidar *Lidar) processClose() (err error) {
	if runtime.GOOS == "windows" {
		log.Println("closing is not implemented on Windows, killing instead")
		return lidar.processKill()
	}

	log.Println("closing lidar-scan process")
	if err = lidar.process.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	lidar.running = false
	return nil
}

// Kills the lidar-scan process, so the cleanup will not be performed.
// Emergancy only.
func (lidar *Lidar) processKill() (err error) {
	log.Println("killing lidar-scan process")
	if err = lidar.process.Process.Kill(); err != nil {
		return err
	}
	lidar.running = false
	return nil
}
