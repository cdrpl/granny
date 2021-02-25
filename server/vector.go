package main

import (
	"math"
)

// Vector is used to represent positions.
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Return the distance from the given vector.
func (v Vector) distance(v2 Vector) float64 {
	x := math.Pow(v.X-v2.X, 2)
	y := math.Pow(v.Y-v2.Y, 2)
	dist := math.Sqrt(x + y)
	return math.Abs(dist)
}

// Add the given vector to the current vector.
func (v *Vector) add(vB Vector) {
	v.X += vB.X
	v.Y += vB.Y
}

// Multiply the vector by the given value.
func (v *Vector) mult(val float64) {
	v.X *= val
	v.Y *= val
}

// Return the magnitude of the vector.
func (v Vector) magnitude() float64 {
	x := v.X * v.X
	y := v.Y * v.Y
	return math.Sqrt(x + y)
}

// Return a directional vector from pos to dest.
func direction(pos, dest Vector) (v Vector) {
	v.X = dest.X - pos.X
	v.Y = dest.Y - pos.Y
	return
}

// Return the normalized vector
func normalize(v Vector) Vector {
	mag := v.magnitude()

	if mag == 0 {
		return Vector{} // zeroed out vector
	}

	v.X /= mag
	v.Y /= mag
	return v
}
