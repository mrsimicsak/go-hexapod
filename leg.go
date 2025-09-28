package main

import (
	"math"

	"gobot.io/x/gobot/drivers/gpio"
)

const (
	LENGTH_TROCHANTER float64 = 50
	LENGTH_FEMUR      float64 = 80
	LENGTH_TIBIA      float64 = 120
	PI                float64 = 3.1415926535897932384626433832795
	HALF_PI           float64 = 1.5707963267948966192313216916398
	TWO_PI            float64 = 6.283185307179586476925286766559
	DEG_TO_RAD        float64 = 0.017453292519943295769236907684886
	RAD_TO_DEG        float64 = 57.295779513082320876798154814105
	EULER             float64 = 2.718281828459045235360287471352
)

var legZeroOffset = Vector3{130, 0, -120}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Leg struct {
	ServoAngles    [3]float32
	Servos         [3]*gpio.ServoDriver
	TargetPosition Vector3
}

func calcLegServoAngles(pos Vector3, mountAngle float64, mirrored bool) (float64, float64, float64) {

	// ###############################################################
	// ### adjust coodinated for loacal space of leg (rotated leg) ###
	// ###############################################################

	var sinAlpha float64 = math.Sin(mountAngle * DEG_TO_RAD)
	var cosAlpha float64 = math.Cos(mountAngle * DEG_TO_RAD)

	var localCoordinates Vector3

	localCoordinates.X = cosAlpha*pos.X + sinAlpha*pos.Y // adjust for leg mounting rotation and zer0 position
	localCoordinates.X += legZeroOffset.X                // adjust for zero position (otherwise the leg origion is (0,0,0))

	localCoordinates.Y = sinAlpha*pos.X - cosAlpha*pos.Y // we have to multiply by -1 to keep a right hand coordinat system
	localCoordinates.Y += legZeroOffset.Y

	localCoordinates.Z = pos.Z + legZeroOffset.Z

	// ###################
	// ### calc angles ###
	// ###################

	var noTrochanter float64 = math.Sqrt(math.Pow(localCoordinates.X, 2)+math.Pow(localCoordinates.Y, 2)) - LENGTH_TROCHANTER // distance between Servo 2 and the tip on the XY Plane
	var servo2TipDistance float64 = math.Sqrt(math.Pow(noTrochanter, 2) + math.Pow(localCoordinates.Z, 2))                    // distance between servo 2 and the tip

	var angle_0 float64 = math.Atan2(localCoordinates.Y, localCoordinates.X)*RAD_TO_DEG + 90

	var angleRightSideTriangle float64 = math.Atan2(localCoordinates.Z, noTrochanter) * RAD_TO_DEG
	var angleUnequalTriangle float64 = math.Acos((math.Pow(LENGTH_TIBIA, 2)-math.Pow(servo2TipDistance, 2)-math.Pow(LENGTH_FEMUR, 2))/(-2*servo2TipDistance*LENGTH_FEMUR)) * RAD_TO_DEG

	var angleUnequalTriangle_2 float64 = math.Acos((math.Pow(servo2TipDistance, 2)-math.Pow(LENGTH_TIBIA, 2)-math.Pow(LENGTH_FEMUR, 2))/(-2*LENGTH_TIBIA*LENGTH_FEMUR)) * RAD_TO_DEG // second angle of unequal triangle

	var angle_1 float64 = 0
	var angle_2 float64 = 0

	if mirrored { // if the leg is mirrored the servos are installed slightly different
		angle_1 = (angleRightSideTriangle+angleUnequalTriangle)*-1 + 90
		angle_2 = (angleUnequalTriangle_2)
	} else {
		angle_1 = (angleRightSideTriangle + angleUnequalTriangle) + 90
		angle_2 = 180 - (angleUnequalTriangle_2)
	}

	// ####################################
	// ### set target angles for servos ###
	// ####################################

	return angle_0, angle_1, angle_2
}
