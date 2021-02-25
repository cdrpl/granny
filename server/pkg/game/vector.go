package game

import (
	"math"
)

// Vector is used to represent positions.
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// MoveTowards will move the vector towards the destination.
func (v *Vector) MoveTowards(dest Vector, speed float64) {
	dir := direction(*v, dest)
	dir = normalize(dir)
	dir = mult(dir, speed)

	dist := distance(*v, dest)

	if dist <= magnitude(dir) {
		*v = dest
	} else {
		v.add(dir)
	}
}

// Return the distance between two vectors.
func distance(v1, v2 Vector) float64 {
	x := math.Pow(v1.X-v2.X, 2)
	y := math.Pow(v1.Y-v2.Y, 2)
	dist := math.Sqrt(x + y)
	return math.Abs(dist)
}

// Add the given vector to the current vector.
func (v *Vector) add(vB Vector) {
	v.X += vB.X
	v.Y += vB.Y
}

// Multiply a vector by the given value
func mult(v Vector, val float64) Vector {
	v.X *= val
	v.Y *= val
	return v
}

// Return the magnitude of the given vector
func magnitude(v Vector) float64 {
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
	mag := magnitude(v)

	if mag == 0 {
		return Vector{} // zeroed out vector
	}

	v.X /= mag
	v.Y /= mag
	return v
}
