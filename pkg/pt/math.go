package pt

import (
	"math"
	"math/rand"
)

type Color Vector3

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func NewVector3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

// Vector with random X,Y,Z components within [min,max]
func NewRandomVector(min, max float64, r *rand.Rand) Vector3 {
	return Vector3{RandFloat(min, max, r), RandFloat(min, max, r), RandFloat(min, max, r)}
}

func (v1 Vector3) Add(v2 Vector3) Vector3 {
	return Vector3{X: v1.X + v2.X, Y: v1.Y + v2.Y, Z: v1.Z + v2.Z}
}

func (v1 Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{X: v1.X - v2.X, Y: v1.Y - v2.Y, Z: v1.Z - v2.Z}
}

func (v1 Vector3) Inverse() Vector3 {
	return Vector3{X: 1 / v1.X, Y: 1 / v1.Y, Z: 1 / v1.Z}
}

func (v1 Vector3) ElemMul(v2 Vector3) Vector3 {
	return Vector3{v1.X * v2.X, v1.Y * v2.Y, v1.Z * v2.Z}
}

func (v1 Vector3) Mul(factor float64) Vector3 {
	return Vector3{X: v1.X * factor, Y: v1.Y * factor, Z: v1.Z * factor}
}

func (v1 Vector3) LengthSquared() float64 {
	return v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z
}

func (v1 Vector3) Length() float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)
}

func (v1 Vector3) Dot(v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector3) Cross(v2 Vector3) Vector3 {
	return Vector3{X: v1.Y*v2.Z - v1.Z*v2.Y, Y: v1.Z*v2.X - v1.X*v2.Z, Z: v1.X*v2.Y - v1.Y*v2.X}
}

func (v1 Vector3) Unit() Vector3 {
	l := v1.Length()
	return Vector3{v1.X / l, v1.Y / l, v1.Z / l}
}

func (v Vector3) ApproxZero() bool {
	// Coparison "hard coded" because it's twice as fast as calling ApproxZero()
	return v.X <= APPROX_THRESH && v.Y <= APPROX_THRESH && v.Z <= APPROX_THRESH
}

// Returns homogenous representation of a point as Vector4
func (p Vector3) ToPoint() Vector4 {
	return Vector4{
		x: p.X,
		y: p.Y,
		z: p.Z,
		w: 1,
	}
}

// Returns homogenous representation of a vector as Vector4
func (p Vector3) ToVector() Vector4 {
	return Vector4{
		x: p.X,
		y: p.Y,
		z: p.Z,
		w: 0,
	}
}

type Vector4 struct {
	x, y, z, w float64
}

func (v Vector4) ToV3() Vector3 {
	return Vector3{
		X: v.x,
		Y: v.y,
		Z: v.z,
	}
}

func (p Vector4) Transformed(m Matrix4) Vector4 {
	return m.Multiply(p)
}

// 4x4 Matrix [y][x], [row][column]
type Matrix4 struct {
	values [4][4]float64
}

func IdentityMatrix() Matrix4 {
	return Matrix4{
		[4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func (m Matrix4) Multiply(v Vector4) Vector4 {
	result := make([]float64, 4)
	for i, row := range m.values {
		sum := row[0] * v.x
		sum += row[1] * v.y
		sum += row[2] * v.z
		sum += row[3] * v.w
		result[i] = sum
	}
	return Vector4{x: result[0], y: result[1], z: result[2], w: result[3]}
}

func (m1 Matrix4) MultiplyMatrix(m2 Matrix4) Matrix4 {
	a := m1.values
	b := m2.values
	return Matrix4{
		[4][4]float64{
			{
				a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0] + a[0][3]*b[3][0],
				a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1] + a[0][3]*b[3][1],
				a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2] + a[0][3]*b[3][2],
				a[0][0]*b[0][3] + a[0][1]*b[1][3] + a[0][2]*b[2][3] + a[0][3]*b[3][3],
			},
			{
				a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0] + a[1][3]*b[3][0],
				a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1] + a[1][3]*b[3][1],
				a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2] + a[1][3]*b[3][2],
				a[1][0]*b[0][3] + a[1][1]*b[1][3] + a[1][2]*b[2][3] + a[1][3]*b[3][3],
			},
			{
				a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0] + a[2][3]*b[3][0],
				a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1] + a[2][3]*b[3][1],
				a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2] + a[2][3]*b[3][2],
				a[2][0]*b[0][3] + a[2][1]*b[1][3] + a[2][2]*b[2][3] + a[2][3]*b[3][3],
			},
			{
				a[3][0]*b[0][0] + a[3][1]*b[1][0] + a[3][2]*b[2][0] + a[3][3]*b[3][0],
				a[3][0]*b[0][1] + a[3][1]*b[1][1] + a[3][2]*b[2][1] + a[3][3]*b[3][1],
				a[3][0]*b[0][2] + a[3][1]*b[1][2] + a[3][2]*b[2][2] + a[3][3]*b[3][2],
				a[3][0]*b[0][3] + a[3][1]*b[1][3] + a[3][2]*b[2][3] + a[3][3]*b[3][3],
			},
		},
	}
}

type Quanternion struct {
	a float64
	v Vector3
}

func NewQuanternion(a float64, v Vector3) Quanternion {
	return Quanternion{
		a: a,
		v: v,
	}
}

func (a Quanternion) Mul(b Quanternion) Quanternion {
	aa := a.a*b.a - a.v.Dot(b.v)
	vv := b.v.Mul(a.a).Add(a.v.Mul(b.a)).Add(a.v.Cross(b.v))
	return NewQuanternion(aa, vv)
}

func (a Quanternion) ToRotationMatrix() Matrix4 {
	w := a.a
	x := a.v.X
	y := a.v.Y
	z := a.v.Z

	vals := [4][4]float64{
		{1 - 2*y*y - 2*z*z, 2*x*y - 2*z*w, 2*x*z + 2*y*w, 0},
		{2*x*y + 2*z*w, 1 - 2*x*x - 2*z*z, 2*y*z - 2*x*w, 0},
		{2*x*z - 2*y*w, 2*y*z + 2*x*w, 1 - 2*x*x - 2*y*y, 0},
		{0, 0, 0, 1},
	}
	return Matrix4{
		values: vals,
	}
}
