package main

// PointSph3 is 3-dimensional point in spherical coordinate system
type PointSph3 struct {
	AngleLidar float32 // Lidar angle in degrees.
	AngleServo float32 // Servo angle in degrees.
	Dist       float32 // Distance in millimeters.
}

// PointCar3 is 3-dimensional point in cartesian coordinate system
type PointCar3 struct {
	X float32
	Y float32
	Z float32
}

func mergerLidarServoV1(lidarBuff *LidarCloud, servoBuff *ServoDataBuffer, accelBuff *AccelDataBuffer) (cloud []PointSph3, err error) {

	return cloud, nil
}
