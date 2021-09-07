package pt

import (
	"math"
	"math/rand"
)

type Color Vector3

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

func NewColor255(r, g, b int) Color {
	return Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
}

func (c1 Color) Blend(c2 Color) Color {
	return Color{c1.X * c2.X, c1.Y * c2.Y, c1.Z * c2.Z}
}

func (c1 Color) Add(c2 Color) Color {
	return Color{c1.X + c2.X, c1.Y + c2.Y, c1.Z + c2.Z}
}

func (c1 Color) Sub(c2 Color) Color {
	return Color{c1.X - c2.X, c1.Y - c2.Y, c1.Z - c2.Z}
}

func (c Color) Div(n float64) Color {
	return Color{c.X / n, c.Y / n, c.Z / n}
}

func (c Color) Scale(n float64) Color {
	return Color{c.X * n, c.Y * n, c.Z * n}
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
	return math.Abs(v.X) <= APPROX_THRESH && math.Abs(v.Y) <= APPROX_THRESH && math.Abs(v.Z) <= APPROX_THRESH
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

// 4x4 Matrix
type Matrix4 [16]float64

func IdentityMatrix() Matrix4 {
	return Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (m Matrix4) Multiply(v Vector4) Vector4 {
	return Vector4{
		x: m[0]*v.x + m[1]*v.x + m[2]*v.x + m[3]*v.x,
		y: m[4]*v.y + m[5]*v.y + m[6]*v.y + m[7]*v.y,
		z: m[8]*v.z + m[9]*v.z + m[10]*v.z + m[11]*v.z,
		w: m[12]*v.w + m[13]*v.w + m[14]*v.w + m[15]*v.w,
	}
}

func (m Matrix4) Inverse() Matrix4 {
	var inv [16]float64

	inv[0] = m[5]*m[10]*m[15] -
		m[5]*m[11]*m[14] -
		m[9]*m[6]*m[15] +
		m[9]*m[7]*m[14] +
		m[13]*m[6]*m[11] -
		m[13]*m[7]*m[10]

	inv[4] = -m[4]*m[10]*m[15] +
		m[4]*m[11]*m[14] +
		m[8]*m[6]*m[15] -
		m[8]*m[7]*m[14] -
		m[12]*m[6]*m[11] +
		m[12]*m[7]*m[10]

	inv[8] = m[4]*m[9]*m[15] -
		m[4]*m[11]*m[13] -
		m[8]*m[5]*m[15] +
		m[8]*m[7]*m[13] +
		m[12]*m[5]*m[11] -
		m[12]*m[7]*m[9]

	inv[12] = -m[4]*m[9]*m[14] +
		m[4]*m[10]*m[13] +
		m[8]*m[5]*m[14] -
		m[8]*m[6]*m[13] -
		m[12]*m[5]*m[10] +
		m[12]*m[6]*m[9]

	inv[1] = -m[1]*m[10]*m[15] +
		m[1]*m[11]*m[14] +
		m[9]*m[2]*m[15] -
		m[9]*m[3]*m[14] -
		m[13]*m[2]*m[11] +
		m[13]*m[3]*m[10]

	inv[5] = m[0]*m[10]*m[15] -
		m[0]*m[11]*m[14] -
		m[8]*m[2]*m[15] +
		m[8]*m[3]*m[14] +
		m[12]*m[2]*m[11] -
		m[12]*m[3]*m[10]

	inv[9] = -m[0]*m[9]*m[15] +
		m[0]*m[11]*m[13] +
		m[8]*m[1]*m[15] -
		m[8]*m[3]*m[13] -
		m[12]*m[1]*m[11] +
		m[12]*m[3]*m[9]

	inv[13] = m[0]*m[9]*m[14] -
		m[0]*m[10]*m[13] -
		m[8]*m[1]*m[14] +
		m[8]*m[2]*m[13] +
		m[12]*m[1]*m[10] -
		m[12]*m[2]*m[9]

	inv[2] = m[1]*m[6]*m[15] -
		m[1]*m[7]*m[14] -
		m[5]*m[2]*m[15] +
		m[5]*m[3]*m[14] +
		m[13]*m[2]*m[7] -
		m[13]*m[3]*m[6]

	inv[6] = -m[0]*m[6]*m[15] +
		m[0]*m[7]*m[14] +
		m[4]*m[2]*m[15] -
		m[4]*m[3]*m[14] -
		m[12]*m[2]*m[7] +
		m[12]*m[3]*m[6]

	inv[10] = m[0]*m[5]*m[15] -
		m[0]*m[7]*m[13] -
		m[4]*m[1]*m[15] +
		m[4]*m[3]*m[13] +
		m[12]*m[1]*m[7] -
		m[12]*m[3]*m[5]

	inv[14] = -m[0]*m[5]*m[14] +
		m[0]*m[6]*m[13] +
		m[4]*m[1]*m[14] -
		m[4]*m[2]*m[13] -
		m[12]*m[1]*m[6] +
		m[12]*m[2]*m[5]

	inv[3] = -m[1]*m[6]*m[11] +
		m[1]*m[7]*m[10] +
		m[5]*m[2]*m[11] -
		m[5]*m[3]*m[10] -
		m[9]*m[2]*m[7] +
		m[9]*m[3]*m[6]

	inv[7] = m[0]*m[6]*m[11] -
		m[0]*m[7]*m[10] -
		m[4]*m[2]*m[11] +
		m[4]*m[3]*m[10] +
		m[8]*m[2]*m[7] -
		m[8]*m[3]*m[6]

	inv[11] = -m[0]*m[5]*m[11] +
		m[0]*m[7]*m[9] +
		m[4]*m[1]*m[11] -
		m[4]*m[3]*m[9] -
		m[8]*m[1]*m[7] +
		m[8]*m[3]*m[5]

	inv[15] = m[0]*m[5]*m[10] -
		m[0]*m[6]*m[9] -
		m[4]*m[1]*m[10] +
		m[4]*m[2]*m[9] +
		m[8]*m[1]*m[6] -
		m[8]*m[2]*m[5]

	det := m[0]*inv[0] + m[1]*inv[4] + m[2]*inv[8] + m[3]*inv[12]

	det = 1.0 / det
	return Matrix4{
		inv[0] * det, inv[1] * det, inv[2] * det, inv[3] * det,
		inv[4] * det, inv[5] * det, inv[6] * det, inv[7] * det,
		inv[8] * det, inv[9] * det, inv[10] * det, inv[11] * det,
		inv[12] * det, inv[13] * det, inv[14] * det, inv[15] * det,
	}
}

func (m Matrix4) Transpose() Matrix4 {
	return Matrix4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

func (a Matrix4) MultiplyMatrix(b Matrix4) Matrix4 {
	return Matrix4{
		a[0]*b[0] + a[1]*b[4] + a[2]*b[8] + a[3]*b[12],
		a[0]*b[1] + a[1]*b[5] + a[2]*b[9] + a[3]*b[13],
		a[0]*b[2] + a[1]*b[6] + a[2]*b[10] + a[3]*b[14],
		a[0]*b[3] + a[1]*b[7] + a[2]*b[11] + a[3]*b[15],

		a[4]*b[0] + a[5]*b[4] + a[6]*b[8] + a[7]*b[12],
		a[4]*b[1] + a[5]*b[5] + a[6]*b[9] + a[7]*b[13],
		a[4]*b[2] + a[5]*b[6] + a[6]*b[10] + a[7]*b[14],
		a[4]*b[3] + a[5]*b[7] + a[6]*b[11] + a[7]*b[15],

		a[8]*b[0] + a[9]*b[4] + a[10]*b[8] + a[11]*b[12],
		a[8]*b[1] + a[9]*b[5] + a[10]*b[9] + a[11]*b[13],
		a[8]*b[2] + a[9]*b[6] + a[10]*b[10] + a[11]*b[14],
		a[8]*b[3] + a[9]*b[7] + a[10]*b[11] + a[11]*b[15],

		a[12]*b[0] + a[13]*b[4] + a[14]*b[8] + a[15]*b[12],
		a[12]*b[1] + a[13]*b[5] + a[14]*b[9] + a[15]*b[13],
		a[12]*b[2] + a[13]*b[6] + a[14]*b[10] + a[15]*b[14],
		a[12]*b[3] + a[13]*b[7] + a[14]*b[11] + a[15]*b[15],
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

	return Matrix4{
		1 - 2*y*y - 2*z*z, 2*x*y - 2*z*w, 2*x*z + 2*y*w, 0,
		2*x*y + 2*z*w, 1 - 2*x*x - 2*z*z, 2*y*z - 2*x*w, 0,
		2*x*z - 2*y*w, 2*y*z + 2*x*w, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}
