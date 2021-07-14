package pt

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
