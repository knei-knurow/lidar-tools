package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// lidar constants
const (
	lidarMaxDataSize = 8192
)

// rplidar scanning modes
const (
	rplidarModeBoost       = 2
	rplidarModeSensitivity = 3
	rplidarModeStability   = 4
	rplidarModeDefault     = rplidarModeSensitivity
)

// Represents a single angle+dist measurement
type Point struct {
	Angle float32 // angle in degrees
	Dist  float32 // distance in millimeters
	// pointpt time.Time // estimated measurement time point
}

// Contains one full-360deg lidar point cloud
type LidarCloud struct {
	Id        int                     //
	TimeBegin time.Time               // time point of starting line read (line starting with '!')
	TimeDiff  int                     // number of milliseconds of current cloud measurement (received from lidar-scan)
	timeEnd   time.Time               // timeBegin increased by timeDiff milliseconds
	Data      [lidarMaxDataSize]Point // measurements data
	Size      uint                    // number of used points in Data
}

// General lidar parameters
type Lidar struct {
	TimeInit time.Time     // Time of the first starting line read (line starting with '!')
	Rpm      int           // declared rpm (but actual may differ)
	Mode     int           // rplidar mode
	Argv     []string      // lidar-scan process argv
	Path     string        // lidar-scan path
	Stdout   io.ReadCloser // lidar-scan stdout
	Stderr   io.ReadCloser // lidar-scan stderr
	process  *exec.Cmd     // lidar-scan process
	running  bool          // whether lidar-scan is currently scanning
}

// Starts the lidar-scan process.
// Does not check whether it has been already started.
func (lidar *Lidar) ProcessStart() error {
	var err error
	lidar.process = exec.Command(lidar.Path, lidar.Argv...)
	log.Println("starting lidar-scan process")

	lidar.Stdout, err = lidar.process.StdoutPipe()
	if err != nil {
		return errors.New("Unable to get stdout of lidar-scan process.")
	}
	lidar.Stderr, err = lidar.process.StderrPipe()
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
func (lidar *Lidar) ProcessClose() (err error) {
	if runtime.GOOS == "windows" {
		log.Println("closing is not implemented on Windows, killing instead")
		return lidar.ProcessKill()
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
func (lidar *Lidar) ProcessKill() (err error) {
	log.Println("killing lidar-scan process")
	if err = lidar.process.Process.Kill(); err != nil {
		return err
	}
	lidar.running = false
	return nil
}

func (lidar *Lidar) LoopStart() (err error) {
	if err != lidar.ProcessStart() {
		return err
	}

	var cloud LidarCloud

	scanner := bufio.NewScanner(lidar.Stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		if err := lidar.LineProcess(line, &cloud); err != nil {
			log.Printf("unable to parse line: %s\n", err)
			// TODO: buffer overflow error handling (but tbh it never happens)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (lidar *Lidar) LineProcess(line string, cloud *LidarCloud) (err error) {
	if len(line) == 0 {
		return
	}

	switch line[0] {
	case '#':
	case '!':
		var cnt, timeDiff int
		if _, err := fmt.Sscanf(line, "! %d %d", &cnt, &timeDiff); err != nil {
			return errors.New("invalid starting line.")
		}

		log.Printf("processed new point cloud (id:%d, timediff:%dms, size:%d)\n", cloud.Id, cloud.TimeDiff, cloud.Size)

		*cloud = LidarCloud{
			Id:        cnt + 1,
			TimeBegin: time.Now(),
			TimeDiff:  timeDiff,
			timeEnd:   time.Now().Add(time.Millisecond * time.Duration(timeDiff)),
		}
	default:
		var angle, dist float32
		if _, err := fmt.Sscanf(line, "%f %f", &angle, &dist); err != nil {
			return fmt.Errorf("invalid data line: \"%s\"\n", line)
		}

		cloud.Size++
		if cloud.Size >= lidarMaxDataSize {
			return errors.New("data buffer overflow")
		}
		cloud.Data[cloud.Size] = Point{angle, dist}
	}
	return nil
}
