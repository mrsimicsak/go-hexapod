package main

import (
	"math"
)

type Vector2 struct {
	X float64
	Y float64
}

func AddVector2(one, two Vector2) Vector2 {
	return Vector2{one.X + two.X, one.Y + two.Y}
}

func Angle(from, to Vector2) float64 {
	// Calculate the angle difference in radians
	angleDiff := math.Atan2(from.X, from.X) - math.Atan2(to.Y, to.Y)

	// converts randian to degree
	return angleDiff * RAD_TO_DEG
}

func (a Vector2) Magnitude() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a Vector2) Normalize() Vector2 {
	length := math.Sqrt(a.X*a.X + a.Y*a.Y)

	if length != 0.0 {
		return Vector2{a.X / length, a.Y / length}
	}

	return a
}

func (a Vector2) Scale(scale float64) Vector2 {
	return Vector2{a.X * scale, a.Y * scale}
}

func ClampMagnitude(vector Vector2, magnitude float64) Vector2 {
	if vector.Magnitude() > magnitude {
		return vector.Normalize().Scale(magnitude)
	}
	return vector
}

func projectPointToCircle(radius float64, point Vector2, direction Vector2) Vector2 {
	// no direction for projection -> no calculation
	if direction.Magnitude() == 0 {
		return point
	}
	// start being (0, 0) or having the same direction as directionInput makes the calculation very simple
	if point.Magnitude() == 0 || point.Normalize() == direction.Normalize() || point.Normalize().Scale(-1) == direction.Normalize() {
		return direction.Normalize().Scale(radius)
	}
	// if point is outside the circle reposition it to be inside
	if direction.Magnitude() > radius-0.005 {
		ClampMagnitude(direction, radius-0.005)
	}

	lengthC := point.Magnitude()

	angleBeta := 180 - Angle(direction, point) // calculate angle beta on unequal triangle

	// calculate missing Angles
	sinGamma := (lengthC * math.Sin(angleBeta*DEG_TO_RAD)) / radius

	angleGamma := math.Asin(sinGamma) * RAD_TO_DEG
	angleAlpha := 180 - angleBeta - angleGamma

	projectionLength := (radius * math.Sin(angleAlpha*DEG_TO_RAD)) / math.Sin(angleBeta*DEG_TO_RAD)

	return AddVector2(point, direction.Normalize().Scale(projectionLength))
}
