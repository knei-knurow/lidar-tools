package main

import (
	"fmt"
	"log"
	"math"
)

// Vec2 represents a (X, Y) vector.
type Vec2 struct {
	X float64 // X
	Y float64 // Y
}

// Vec3 represents a (X, Y, Z) vector.
type Vec3 struct {
	X float64 // X
	Y float64 // Y
	Z float64 // Z
}

// Quat represents a quaternion.
type Quat struct {
	W float64
	X float64
	Y float64
	Z float64
}

// RadToDeg converts radians to degrees.
func RadToDeg(a float64) float64 {
	return a * 180 / math.Pi
}

// DegToRad converts degrees to radians.
func DegToRad(a float64) float64 {
	return a * math.Pi / 180
}

// AngleDistToPoint2 converts angle + dist pair to (X, Y) vector.
func AngleDistToPoint2(v *AngleDist) (w Vec2) {
	w.X = v.Dist * math.Cos(DegToRad(v.Angle))
	w.Y = v.Dist * math.Sin(DegToRad(v.Angle))
	return w
}

// QuatMult multiplies two quaternions.
func QuatMult(q1 *Quat, q2 *Quat) (q3 Quat) {
	w1, x1, y1, z1 := q1.W, q1.X, q1.Y, q1.Z
	w2, x2, y2, z2 := q2.W, q2.X, q2.Y, q2.Z
	q3.W = w1*w2 - x1*x2 - y1*y2 - z1*z2
	q3.X = w1*x2 + x1*w2 + y1*z2 - z1*y2
	q3.Y = w1*y2 + y1*w2 + z1*x2 - x1*z2
	q3.Z = w1*z2 + z1*w2 + x1*y2 - y1*x2
	return
}

// QuatConjugate returns quaternion conjugate.
func QuatConjugate(q *Quat) Quat {
	return Quat{q.W, -q.X, -q.Y, -q.Z}
}

// QuatVec3Mult performs quaternion-vector multiplication.
func QuatVec3Mult(q1 *Quat, v *Vec3) Vec3 {
	q2 := Quat{0, v.X, v.Y, v.Z}

	a := QuatMult(q1, &q2)
	b := QuatConjugate(q1)
	w := QuatMult(&a, &b)

	return Vec3{w.X, w.Y, w.Z}
}

// RotateVec3ByQuat rotates (x, y, z) vector by a normalised quaternion.
func RotateVec3ByQuat(v *Vec3, q *Quat) (w Vec3) {
	return QuatVec3Mult(q, v)
}

// RotateVec2 rotates (x, y) around origin (0, 0) by a rads.
func RotateVec2(v *Vec2, a float64) (w Vec2) {
	w.X = v.X*math.Cos(a) - v.Y*math.Sin(a)
	w.Y = v.Y*math.Cos(a) + v.X*math.Sin(a)
	return
}

type Fusion struct {
	CloudRotation float64 // each scanned 2D cloud will be rotated by CloudRotation radians
	cloudsCnt     uint
}

const (
	PrototypeCloudRotation = -math.Pi / 4 // cloud rotation for the first lidar head prototype
)

func (fusion *Fusion) UpdateWithAccel(cloud *LidarCloud, accel *AccelDataBuffer) {
	if cloud.Size == 0 {
		return
	}

	// estimated time of a single point measurement
	// timePerPt := float64(cloud.TimeDiff) / float64(cloud.Size) * 1000

	var q0 AccelDataQuat
	for j := 0; j < accel.size; j++ { // from the latest to the earliest
		a0, _ := accel.Get(j)
		if a0.quat.timept.Before(cloud.TimeBegin) {
			q0 = a0.quat
			break
		}
	}

	for i := 0; i < int(cloud.Size); i++ {
		if cloud.Data[i].Dist == 0 {
			continue
		}

		// estimated time point of the i-th measurement
		// t := cloud.TimeBegin.Add(time.Microsecond * time.Duration(timePerPt*float64(i)))

		// POSSIBLE ERROR SOURCE: there was an idea to take an average of two accel measurements
		// which would be biased towards the later or earlier one (depending on the t value).
		// This approach requires additional quaternion computation, more info here:
		// https://math.stackexchange.com/q/162863/527542

		// 1. convert (angle, dist) to (X, Y)
		pt2 := AngleDistToPoint2(&cloud.Data[i])

		// 2. rotate lidar cloud depending on the head construction
		pt2 = RotateVec2(&pt2, fusion.CloudRotation)

		// 3. modify (X, Y) to (X, Y, Z) where Z=0
		pt3 := Vec3{pt2.X, pt2.Y, 0}

		// 4. rotate (X, Y, Z) by accel quaternion to get (X', Y', Z')
		pt3 = RotateVec3ByQuat(&pt3, &Quat{q0.qw, q0.qx, q0.qy, q0.qz})

		fmt.Printf("%f\t%f\t%f\n", pt3.X, pt3.Y, pt3.Z)
	}

	fusion.cloudsCnt++
}

func (fusion *Fusion) UpdateWithServo(cloud *LidarCloud, servoData *ServoDataBuffer, servo *Servo) {
	if cloud.Size == 0 {
		return
	}

	var s0 ServoData
	for j := 0; j < servoData.size; j++ { // from the latest to the earliest
		s0, _ = servoData.Get(j)
		if s0.timept.Before(cloud.TimeBegin) {
			break
		}
	}

	for i := 0; i < int(cloud.Size); i++ {
		if cloud.Data[i].Dist == 0 {
			// continue
		}

		// 1. convert (angle, dist) to (X, Y)
		pt2 := AngleDistToPoint2(&cloud.Data[i])

		// 2. rotate lidar cloud depending on the head construction
		pt2 = RotateVec2(&pt2, fusion.CloudRotation)

		// // 3.
		deg := (float64(s0.positon) - float64(servo.positonStart)) * servoUnitToDeg
		pt2t := RotateVec2(&Vec2{pt2.X, 0}, DegToRad(float64(deg)))

		// 4.
		pt3 := Vec3{pt2t.X, pt2.Y, pt2t.Y}
		fmt.Printf("%f\t%f\t%f\n", pt3.X, pt3.Y, pt3.Z)
	}
	log.Println((float64(s0.positon) - float64(servo.positonStart)) * servoUnitToDeg)
	fusion.cloudsCnt++
}
