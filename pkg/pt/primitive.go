package pt

import (
	"math"
)

type hit struct {
	point     Vector3  // intersection Point
	normal    Vector3  // normal at the intersection Point always pointing agains the ray
	frontFace bool     // Wheter or not the ray hit from the outside or the inside
	t         float64  // distance along the intersection ray
	material  Material // Material at intersection point
}

type intersectable interface {
	intersected(ray ray, tMin, tMax float64, hitOut *hit) bool
	bounding() aabb
}

type primitive interface {
	intersectable
	transformed(Matrix4) primitive
}

type Sphere struct {
	center Vector3
	radius float64
	box    aabb
}

func NewSphere(center Vector3, radius float64) *Sphere {
	radVec := NewVector3(radius, radius, radius)
	min := center.Sub(radVec)
	max := center.Add(radVec)
	return &Sphere{
		center: center,
		radius: radius,
		box:    newAABB(min, max),
	}
}

func (s *Sphere) transformed(t Matrix4) primitive {
	/*
		pointOnSphere := s.center.Add(NewVector3(s.radius, 0, 0))
		newCenter := s.center.ToPoint().Transformed(t).ToV3()
		newRadius := pointOnSphere.ToPoint().Transformed(t).ToV3().Length()
		return NewSphere(newCenter, newRadius)
	*/
	return NewSphere(s.center.ToPoint().Transformed(t).ToV3(), s.radius)
}

func (s *Sphere) bounding() aabb {
	return s.box
}

func (s *Sphere) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	oc := ray.origin.Sub(s.center)
	dirNorm := ray.direction.Length()
	a := dirNorm * dirNorm
	halfB := oc.Dot(ray.direction)
	ocNorm := oc.Length()
	c := ocNorm*ocNorm - s.radius*s.radius
	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false
	}

	// Nearest intersection distance within tMin <= t <= tMax
	sqrtDiscriminant := math.Sqrt(discriminant)
	t := (-halfB - sqrtDiscriminant) / a
	if t <= tMin || t >= tMax {
		t = (-halfB + sqrtDiscriminant) / a
		if t <= tMin || t >= tMax {
			return false
		}
	}

	hitOut.point = ray.position(t)
	hitOut.normal = hitOut.point.Sub(s.center).Mul(1 / s.radius)
	hitOut.frontFace = ray.direction.Dot(hitOut.normal) < 0
	if !hitOut.frontFace {
		hitOut.normal = hitOut.normal.Mul(-1)
	}
	hitOut.t = t
	return true
}

type vertex struct {
	position Vector3
	normal   Vector3
}

type Triangle struct {
	vertecies [3]vertex
	box       aabb

	// Precalculate v0v1 and v0v2 as it's used
	v0v1 Vector3
	v0v2 Vector3
}

func NewTriangle(vertecies [3]vertex) *Triangle {
	x := [3]float64{vertecies[0].position.X, vertecies[1].position.X, vertecies[2].position.X}
	y := [3]float64{vertecies[0].position.Y, vertecies[1].position.Y, vertecies[2].position.Y}
	z := [3]float64{vertecies[0].position.Z, vertecies[1].position.Z, vertecies[2].position.Z}
	min := NewVector3(Min3(x), Min3(y), Min3(z))
	max := NewVector3(Max3(x), Max3(y), Max3(z))
	return &Triangle{
		box:       newAABB(min, max),
		vertecies: vertecies,
		v0v1:      vertecies[1].position.Sub(vertecies[0].position),
		v0v2:      vertecies[2].position.Sub(vertecies[0].position),
	}
}

func NewTriangleWithoutNormals(v0 Vector3, v1 Vector3, v2 Vector3) *Triangle {
	vertecies := [3]vertex{
		{
			position: v0,
			normal:   calcNormal(v0, v1, v2),
		},
		{
			position: v1,
			normal:   calcNormal(v1, v2, v0),
		},
		{
			position: v2,
			normal:   calcNormal(v2, v0, v1),
		},
	}
	return NewTriangle(vertecies)
}

func calcNormal(point Vector3, right Vector3, left Vector3) Vector3 {
	pa := left.Sub(point)
	pb := right.Sub(point)
	return pb.Cross(pa).Unit()
}

func (t *Triangle) bounding() aabb {
	return t.box
}

// Takes u and v barycentric coordinates and returns the normal at point p
func (tri *Triangle) normal(u, v float64) Vector3 {
	normalW := tri.vertecies[0].normal.Mul(1 - u - v)
	normalU := tri.vertecies[1].normal.Mul(u)
	normalV := tri.vertecies[2].normal.Mul(v)
	return normalU.Add(normalV).Add(normalW)
}

func (tri *Triangle) transformed(t Matrix4) primitive {

	tinv := t.Transpose().Inverse()

	var vertecies [3]vertex
	vertecies[0] = vertex{
		position: tri.vertecies[0].position.ToPoint().Transformed(t).ToV3(),
		normal:   tri.vertecies[0].normal.ToPoint().Transformed(tinv).ToV3(),
	}
	vertecies[1] = vertex{
		position: tri.vertecies[1].position.ToPoint().Transformed(t).ToV3(),
		normal:   tri.vertecies[1].normal.ToPoint().Transformed(tinv).ToV3(),
	}
	vertecies[2] = vertex{
		position: tri.vertecies[2].position.ToPoint().Transformed(t).ToV3(),
		normal:   tri.vertecies[2].normal.ToPoint().Transformed(tinv).ToV3(),
	}
	return NewTriangle(vertecies)
}

func (tri *Triangle) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	// Implementation of the Möller-Trumbore algorithm
	pvec := ray.direction.Cross(tri.v0v2)
	det := tri.v0v1.Dot(pvec)

	// If det is close to 0, Triangle and ray are parallel => no intersection
	if ApproxZero(det) {
		return false
	}

	invDet := 1 / det
	tvec := ray.origin.Sub(tri.vertecies[0].position)
	u := tvec.Dot(pvec) * invDet
	if u < 0 || u > 1 {
		return false
	}

	qvec := tvec.Cross(tri.v0v1)
	v := ray.direction.Dot(qvec) * invDet
	if v < 0 || u+v > 1 {
		return false
	}

	t := tri.v0v2.Dot(qvec) * invDet
	if t < tMin || t > tMax {
		return false
	}

	hitOut.point = ray.position(t)
	hitOut.frontFace = det > 0
	hitOut.normal = tri.normal(u, v)
	if !hitOut.frontFace {
		hitOut.normal = hitOut.normal.Mul(-1)
	}
	hitOut.t = t
	return true
}

// Wrapper for pimitive including a material
type tracable struct {
	prim primitive
	mat  Material
}

func (p tracable) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	if p.prim.intersected(ray, tMin, tMax, hitOut) {
		hitOut.material = p.mat
		return true
	}
	return false
}

func (p tracable) bounding() aabb {
	return p.prim.bounding()
}
