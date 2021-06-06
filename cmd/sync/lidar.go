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

// Lidar-related constants.
const (
	lidarMaxDataSize = 8192

	// Lidar scanning modes.
	rplidarModeBoost       = 2
	rplidarModeSensitivity = 3
	rplidarModeStability   = 4
	rplidarModeDefault     = rplidarModeSensitivity
)

// Point is a single angle + distance measurement.
type Point struct {
	Angle float32 // Angle in degrees.
	Dist  float32 // Distance in millimeters.
}

// LidarCloud is one full (360deg) lidar point cloud.
type LidarCloud struct {
	ID        int                     //
	TimeBegin time.Time               // time point of starting line read (line starting with '!').
	TimeDiff  int                     // number of milliseconds of current cloud measurement (received from lidar-scan).
	timeEnd   time.Time               // timeBegin increased by timeDiff milliseconds
	Data      [lidarMaxDataSize]Point // Measurements data.
	Size      uint                    // Number of used points in Data.
	Ready     bool
}

// Lidar represents general lidar parameters.
type Lidar struct {
	TimeInit           time.Time     // Time of the first starting line read (line starting with '!').
	RPM                int           // Declared RPM (actual may differ).
	Mode               int           // rplidar scan mode.
	Args               []string      // lidar-scan process argv.
	Path               string        // Path to lidar-scan executable.
	Stdout             io.ReadCloser // lidar-scan stdout.
	Stderr             io.ReadCloser // lidar-scan stderr.
	process            *exec.Cmd     // lidar-scan process.
	running            bool          // Whether lidar-scan is currently scanning.
	nextCloudCount     int
	nextCloudTimeDiff  int
	nextCloudTimeBegin time.Time
}

// StartProcess starts the lidar-scan process.
// It does not check whether it has been already started.
func (lidar *Lidar) StartProcess() error {
	lidar.process = exec.Command(lidar.Path, lidar.Args...)
	log.Println("starting lidar-scan process with args:", lidar.Args)

	var err error
	lidar.Stdout, err = lidar.process.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get stdout of lidar-scan process: %v", err)
	}
	lidar.process.Stderr = os.Stderr
	lidar.Stderr = os.Stderr

	err = lidar.process.Start()
	if err != nil {
		return fmt.Errorf("start lidar-scan process: %v", err)
	}

	lidar.running = true
	return nil
}

// CloseProcess sends SIGINT (ctrl+c) to the lidar-scan process.
// It is important because lidar-scan performs cleanup on SIGINT.
// It does not check whether it has been already started.
//
// Does not work on Windows - instead, it just calls KillProcess.
func (lidar *Lidar) CloseProcess() (err error) {
	if runtime.GOOS == "windows" {
		log.Println("closing is not implemented on Windows, killing instead")
		return lidar.KillProcess()
	}

	log.Println("closing lidar-scan process")
	if err = lidar.process.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("send SIGINT to lidar-scan: %v", err)
	}
	lidar.running = false
	return nil
}

// KillProcess kills the lidar-scan process immediately, so the cleanup will not be performed.
//
// Use it only in emergency situations. Prefer CloseProcess.
func (lidar *Lidar) KillProcess() (err error) {
	log.Println("killing lidar-scan process")
	if err = lidar.process.Process.Kill(); err != nil {
		return fmt.Errorf("kill lidar-scan: %v", err)
	}
	lidar.running = false
	return nil
}

// StartLoop starts the lidar-scan process and runs a loop responsible for reading and
// processing lidar data from redirected lidar-scan's stdout. It is designed to be run in a
// goroutine. The channel sends pointers to LidarCloud which contains the latest scanned
// point cloud. The pointers approach is required because LidarCloud is greather than 64kB
// which is a Go limit.
func (lidar *Lidar) StartLoop(channel chan *LidarCloud) (err error) {
	if err := lidar.StartProcess(); err != nil {
		return fmt.Errorf("start process: %v", err)
	}

	scanner := bufio.NewScanner(lidar.Stdout)
	scanner.Split(bufio.ScanLines)
	for {
		// create new cloud every time to pass the pointer via channel and avoid data race
		cloud := LidarCloud{
			ID:        lidar.nextCloudCount + 1,
			TimeDiff:  lidar.nextCloudTimeDiff,
			TimeBegin: lidar.nextCloudTimeBegin, // POSSIBLE ERROR SOURCE: using milliseconds by lidar-scan
			timeEnd:   lidar.nextCloudTimeBegin.Add(time.Millisecond * time.Duration(lidar.nextCloudTimeDiff))}

		for scanner.Scan() {
			line := scanner.Text()

			if err := lidar.ProcessLine(line, &cloud); err != nil {
				log.Printf("unable to parse line: %s\n", err)
				// TODO: buffer overflow error handling (but tbh it never happens)
			}

			if cloud.Ready {
				channel <- &cloud
				break // in order to create new LidarCloud
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
}

// ProcessLine takes a single line from lidar-scan stdout, processes it, and modifies cloud.
func (lidar *Lidar) ProcessLine(line string, cloud *LidarCloud) (err error) {
	if len(line) == 0 {
		return
	}

	switch line[0] {
	case '#':
	case '!':
		lidar.nextCloudTimeBegin = time.Now() // POSSIBLE ERROR SOURCE: using time of data receive
		if _, err := fmt.Sscanf(line, "! %d %d", &lidar.nextCloudCount, &lidar.nextCloudTimeDiff); err != nil {
			return errors.New("invalid starting line")
		}
		cloud.Ready = true
	default:
		var angle, dist float32
		if _, err := fmt.Sscanf(line, "%f %f", &angle, &dist); err != nil {
			return fmt.Errorf("invalid data line: \"%s\"", line)
		}

		cloud.Size++
		if cloud.Size >= lidarMaxDataSize {
			return errors.New("data buffer overflow")
		}
		cloud.Data[cloud.Size] = Point{angle, dist}
	}
	return nil
}

// UpdateProcessArgs modifies the parametres passed to lidar-scan to change its behavior.
//
// In practice, there is no way to modify the command-line arguments passed to lidar-scan
// while it is running, so this function simply kills the process and starts it again
// with updated args.
// Such a solution should be sufficient for most cases.
func (lidar *Lidar) UpdateProcessArgs(args []string) (err error) {
	log.Println("changing lidar-scan process args")
	if err := lidar.CloseProcess(); err != nil {
		return err
	}

	lidar.Args = args

	if err := lidar.StartProcess(); err != nil {
		return err
	}
	return nil
}
