package utils

import (
	"math"
)

var toRad float64 = math.Pi / 180.0

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x float64, y float64) Vec2 {
	return Vec2{X: x, Y: y}
}

func ZeroVec2() Vec2 {
	return NewVec2(0, 0)
}

func float64Equals(a, b float64) bool {
	if a == b {
		return true
	}

	epsilon := 1e-9

	diff := math.Abs(a - b)

	return diff < epsilon

	// if a == 0.0 || b == 0.0 || diff < math.SmallestNonzeroFloat64 {
	// 	return diff < epsilon*math.SmallestNonzeroFloat64
	// }

	// return diff/(math.Abs(a)+math.Abs(b)) < epsilon
}

func (v *Vec2) Set(x float64, y float64) *Vec2 {
	v.X = x
	v.Y = y
	return v
}

func (v Vec2) Equals(other Vec2) bool {
	isXNaN := math.IsNaN(v.X) && math.IsNaN(other.X)
	isYNaN := math.IsNaN(v.Y) && math.IsNaN(other.Y)
	return (isXNaN || float64Equals(v.X, other.X)) && (isYNaN || float64Equals(v.Y, other.Y))
}

func (v Vec2) Add(other Vec2) Vec2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v Vec2) Sub(other Vec2) Vec2 {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v Vec2) Perp() Vec2 {
	return NewVec2(-v.Y, v.X)
}

func (v Vec2) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.LengthSq())
}

func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l > 0 {
		invLength := 1 / v.Length()
		v.X *= invLength
		v.Y *= invLength
	}
	return v
}

func (v Vec2) Rotate(origin Vec2, angleInDeg float64) Vec2 {
	angleInRad := angleInDeg * toRad
	sin, cos := math.Sincos(angleInRad)
	return NewVec2(
		(v.X-origin.X)*cos-(v.Y-origin.Y)*sin,
		(v.X-origin.X)*sin+(v.Y-origin.Y)*cos,
	).Add(origin)
}

func (v Vec2) Scale(n float64) Vec2 {
	v.X *= n
	v.Y *= n
	return v
}

func (v Vec2) ScaleX(n float64) Vec2 {
	v.X *= n
	return v
}

func (v Vec2) ScaleY(n float64) Vec2 {
	v.Y *= n
	return v
}

func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Project on a direction unit vector
func (v Vec2) ProjectNormal(dir Vec2) Vec2 {
	m := v.Dot(dir)
	return dir.Scale(m)
}

func (v Vec2) Project(other Vec2) Vec2 {
	m := v.Dot(other) / other.Length()
	return other.Unit().Scale(m)
}

func (v Vec2) Reflect(normal Vec2) Vec2 {
	return v.Sub(v.ProjectNormal(normal).Scale(2))
}

func (v Vec2) To(other Vec2) Vec2 {
	return other.Sub(v)
}

func (v Vec2) Unit() Vec2 {
	return v.Scale(1 / v.Length())
}

func (v *Vec2) Floor() *Vec2 {
	v.X = math.Floor(v.X)
	v.Y = math.Floor(v.Y)
	return v
}

func (v Vec2) Floored() Vec2 {
	return *v.Floor()
}

func (v Vec2) Cross(other Vec2) float64 {
	return v.X*other.Y - v.Y*other.X
}

func (v Vec2) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vec2) IsNaN() bool {
	return math.IsNaN(v.X) && math.IsNaN(v.Y)
}
