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

type Intersectable interface {
	intersected(ray ray, tMin, tMax float64, hitOut *hit) bool
	bounding() aabb
}

type Primitive interface {
	Intersectable
	transformed(t Matrix4) Primitive
}

// TODO: Is there a better way, than to add mat to primitives?
type sphere struct {
	center Vector3
	radius float64
	mat    Material
}

func newSphere(center Vector3, radius float64, material Material) *sphere {
	return &sphere{
		center: center,
		radius: radius,
		mat:    material,
	}
}

func (s *sphere) transformed(t Matrix4) Primitive {
	// TODO: Implement Me!
	return s
}

func (s *sphere) bounding() aabb {
	// TODO: Cache!
	radVec := NewVector3(s.radius, s.radius, s.radius)
	min := s.center.Sub(radVec)
	max := s.center.Add(radVec)
	return newAABB(min, max)
}

func (s *sphere) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
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

	hitOut.point = ray.Position(t)
	hitOut.normal = hitOut.point.Sub(s.center).Mul(1 / s.radius)
	hitOut.frontFace = ray.direction.Dot(hitOut.normal) < 0
	if !hitOut.frontFace {
		hitOut.normal = hitOut.normal.Mul(-1)
	}
	hitOut.t = t
	hitOut.material = s.mat
	return true
}

type vertex struct {
	position Vector3
	normal   Vector3
}

type triangle struct {
	vertecies [3]vertex
	mat       Material

	// Precalculate v0v1 and v0v2 as it's used
	v0v1 Vector3
	v0v2 Vector3
}

func newTriangle(vertecies [3]vertex, material Material) *triangle {
	return &triangle{
		vertecies: vertecies,
		mat:       material,
		v0v1:      vertecies[1].position.Sub(vertecies[0].position),
		v0v2:      vertecies[2].position.Sub(vertecies[0].position),
	}
}

func newTriangleWithoutNormals(v0 Vector3, v1 Vector3, v2 Vector3, material Material) *triangle {
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
	return newTriangle(vertecies, material)
}

func calcNormal(point Vector3, right Vector3, left Vector3) Vector3 {
	pa := left.Sub(point)
	pb := right.Sub(point)
	return pb.Cross(pa).Unit()
}

func (t *triangle) bounding() aabb {
	// TODO: Cache!
	x := [3]float64{t.vertecies[0].position.X, t.vertecies[1].position.X, t.vertecies[2].position.X}
	y := [3]float64{t.vertecies[0].position.Y, t.vertecies[1].position.Y, t.vertecies[2].position.Y}
	z := [3]float64{t.vertecies[0].position.Z, t.vertecies[1].position.Z, t.vertecies[2].position.Z}
	min := NewVector3(Min3(x), Min3(y), Min3(z))
	max := NewVector3(Max3(x), Max3(y), Max3(z))
	return newAABB(min, max)
}

// Takes u and v barycentric coordinates and returns the normal at point p
func (tri *triangle) normal(u, v float64) Vector3 {
	normalW := tri.vertecies[0].normal.Mul(1 - u - v)
	normalU := tri.vertecies[1].normal.Mul(u)
	normalV := tri.vertecies[2].normal.Mul(v)
	return normalU.Add(normalV).Add(normalW)
}

func (tri *triangle) transformed(t Matrix4) Primitive {
	// TODO: Implement Me!
	return tri
}

func (tri *triangle) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
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

	hitOut.point = ray.Position(t)
	hitOut.frontFace = det > 0
	hitOut.normal = tri.normal(u, v)
	if !hitOut.frontFace {
		hitOut.normal = hitOut.normal.Mul(-1)
	}
	hitOut.t = t
	hitOut.material = tri.mat
	return true
}
