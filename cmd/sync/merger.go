package main

// import (
// 	"fmt"
// 	"log"
// 	"time"
// )

// // PointSph3 is 3-dimensional point in spherical coordinate system
// type PointSph3 struct {
// 	AngleLidar float32 // Lidar angle in degrees.
// 	AngleServo float32 // Servo angle in degrees.
// 	Dist       float32 // Distance in millimeters.
// }

// // PointCar3 is 3-dimensional point in cartesian coordinate system
// type PointCar3 struct {
// 	X float32
// 	Y float32
// 	Z float32
// }

// func buffersPrint(lidarBuff *LidarCloud, servoBuff *ServoDataBuffer, accelBuff *AccelDataBuffer) {
// 	log.Println("-------------------------------")

// 	if lidarBuff != nil {
// 		log.Printf("lidar buffer:\n\tid:\t%d\n\tsize:\t%d\n\tstart:\t%d\n\tend:\t%d\n",
// 			lidarBuff.ID, lidarBuff.Size, lidarBuff.TimeBegin.UnixNano(), lidarBuff.timeEnd.UnixNano())
// 	}

// 	log.Printf("servo buffer:\n")
// 	for i := 0; i < servoBuff.size; i++ {
// 		v, _ := servoBuff.Get(i)
// 		log.Printf("\tpos:\t%d\n", v.positon)
// 		log.Printf("\ttime:\t%d\n", v.timept.UnixNano())
// 	}

// 	log.Printf("accel buffer:\n")
// 	for i := 0; i < accelBuff.size; i++ {
// 		v, _ := accelBuff.Get(i)
// 		log.Printf("\taccel:\t%d\t%d\t%d\n", v.xAccel, v.yAccel, v.zAccel)
// 		log.Printf("\tgyro:\t%d\t%d\t%d\n", v.xGyro, v.yGyro, v.zGyro)
// 		log.Printf("\ttime:\t%d\n", v.timept.UnixNano())
// 	}

// }

// func mergerLidarServoV1(lidarBuff *LidarCloud, servoBuff *ServoDataBuffer, print bool) (cloud []PointSph3) {
// 	if lidarBuff == nil {
// 		return cloud
// 	}

// 	s, _ := servoBuff.Get(0)
// 	servoAngle := float32(s.positon)

// 	cloud = make([]PointSph3, lidarBuff.Size)
// 	for i := 0; i < int(lidarBuff.Size); i++ {
// 		cloud[i] = PointSph3{
// 			lidarBuff.Data[i].Angle,
// 			servoAngle,
// 			lidarBuff.Data[i].Dist,
// 		}

// 		if print && lidarBuff.Data[i].Dist != 0 {
// 			fmt.Println(lidarBuff.Data[i].Angle, servoAngle, lidarBuff.Data[i].Dist)
// 		}
// 	}

// 	return cloud
// }

// func isTimeBetween(v time.Time, start time.Time, end time.Time) bool {
// 	if (v.After(start) && v.Before(end)) || v == start || v == end {
// 		return true
// 	}
// 	return false
// }
