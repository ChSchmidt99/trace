package pt

import (
	"math"
	"math/rand"
)

const APPROX_THRESH = 1e-7

func ApproxZero(num float64) bool {
	return num <= APPROX_THRESH
}

// Deprecated, use RandomUnitVector instead
func RandomUnitVectorOld(r *rand.Rand) Vector3 {
	// First pick a random point in a unit cube, and skip it, if it's outside of the unit sphere
	for {
		v := NewRandomVector(-1, 1, r)
		if v.LengthSquared() >= 1 {
			continue
		}
		return v.Unit()
	}
}

// Generate a random unit vector within a unit sphere
func RandomUnitVector(r *rand.Rand) Vector3 {
	u := r.Float64()
	x1 := r.NormFloat64()
	x2 := r.NormFloat64()
	x3 := r.NormFloat64()
	mag := math.Sqrt(x1*x1 + x2*x2 + x3*x3)
	x1 /= mag
	x2 /= mag
	x3 /= mag
	c := math.Cbrt(u)
	return NewVector3(x1*c, x2*c, x3*c).Unit()
}

/*
func RandomInUnitDisk(r *rand.Rand) Vector3 {
	for {
		v := NewVector3(RandFloat(-1, 1, r), RandFloat(-1, 1, r), 0)
		if v.LengthSquared() >= 1 {
			continue
		}
		return v
	}
}
*/

func RandFloat(min, max float64, r *rand.Rand) float64 {
	return min + r.Float64()*(max-min)
}

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func DegreesToRadians(degree float64) float64 {
	return degree * (math.Pi / 180)
}

// Efficient min of 3 values
func Min3(vals [3]float64) float64 {
	if vals[0] <= vals[1] && vals[0] <= vals[2] {
		return vals[0]
	}
	if vals[1] <= vals[0] && vals[1] <= vals[2] {
		return vals[1]
	}
	return vals[2]
}

// Efficient max of 3 values
func Max3(vals [3]float64) float64 {
	if vals[0] >= vals[1] && vals[0] >= vals[2] {
		return vals[0]
	}
	if vals[1] >= vals[0] && vals[1] >= vals[2] {
		return vals[1]
	}
	return vals[2]
}

// Calculates a vector with the minimum components of a and b combined
func MinVec(a Vector3, b Vector3) Vector3 {
	x := a.X
	if a.X > b.X {
		x = b.X
	}
	y := a.Y
	if a.Y > b.Y {
		y = b.Y
	}
	z := a.Z
	if a.Z > b.Z {
		z = b.Z
	}
	return NewVector3(x, y, z)
}

// Calculates a vector with the maximum components of a and b combined
func MaxVec(a Vector3, b Vector3) Vector3 {
	x := a.X
	if a.X < b.X {
		x = b.X
	}
	y := a.Y
	if a.Y < b.Y {
		y = b.Y
	}
	z := a.Z
	if a.Z < b.Z {
		z = b.Z
	}
	return NewVector3(x, y, z)
}
