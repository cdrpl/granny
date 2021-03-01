package main

import (
	"math"
)

// Vector is used to represent positions.
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Return the magnitude of the vector.
func (v Vector) magnitude() float64 {
	x := v.X * v.X
	y := v.Y * v.Y
	return math.Sqrt(x + y)
}

// Return the normalized vector
func (v Vector) normalize() Vector {
	mag := v.magnitude()

	if mag == 0 {
		return Vector{} // zeroed out vector
	}

	return Vector{
		X: v.X / mag,
		Y: v.Y / mag,
	}
}

// Add the two vectors together.
func vectorAdd(v1, v2 Vector) Vector {
	return Vector{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

// Multiply the vector by the given value.
func vectorMult(v Vector, val float64) Vector {
	return Vector{
		X: v.X * val,
		Y: v.Y * val,
	}
}

// Return the distance from the given vector.
func vectorDist(v1, v2 Vector) float64 {
	x := math.Pow(v1.X-v2.X, 2)
	y := math.Pow(v1.Y-v2.Y, 2)
	dist := math.Sqrt(x + y)
	return math.Abs(dist)
}

// Return a directional vector from pos to dest.
func vectorDir(src, dest Vector) Vector {
	return Vector{
		X: dest.X - src.X,
		Y: dest.Y - src.Y,
	}
}
