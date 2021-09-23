package pt

import (
	"math"
)

type CameraTransformation struct {
	LookFrom Vector3
	LookAt   Vector3
	Up       Vector3
}

type orientation struct {
	origin Vector3
	up     Vector3
	w      Vector3
	u      Vector3
	v      Vector3
}

func newOrientation(t CameraTransformation) orientation {
	w := t.LookFrom.Sub(t.LookAt).Unit()
	u := t.Up.Cross(w).Unit()
	v := w.Cross(u)
	return orientation{
		origin: t.LookFrom,
		w:      t.LookFrom.Sub(t.LookAt).Unit(),
		up:     t.Up,
		u:      u,
		v:      v,
	}
}

type Camera struct {
	orientation orientation

	viewportWidth  float64
	viewportHeight float64

	lowerLeftCorner Vector3
	horizontal      Vector3
	vertical        Vector3
}

func NewDefaultCamera(aspectRatio float64, fov float64) *Camera {
	return NewCamera(aspectRatio, fov, CameraTransformation{
		LookFrom: NewVector3(0, 0, 0),
		LookAt:   NewVector3(1, 0, 0),
		Up:       NewVector3(0, 1, 0),
	})
}

func NewCamera(aspectRatio float64, fov float64, transform CameraTransformation) *Camera {
	cam := new(Camera)
	cam.viewportHeight = 2.0 * math.Tan(DegreesToRadians(fov)/2)
	cam.viewportWidth = aspectRatio * cam.viewportHeight
	cam.SetTransformation(transform)
	return cam
}

func (c *Camera) SetTransformation(transformation CameraTransformation) {
	c.orientation = newOrientation(transformation)
	c.horizontal = c.orientation.u.Mul(c.viewportWidth)
	c.vertical = c.orientation.v.Mul(c.viewportHeight)
	c.lowerLeftCorner = c.orientation.origin.Sub(c.horizontal.Mul(0.5)).Sub(c.vertical.Mul(0.5)).Sub(c.orientation.w)
}

func (c *Camera) Translate(v Vector3) {
	c.orientation.origin = c.orientation.origin.Add(v)
	c.lowerLeftCorner = c.orientation.origin.Sub(c.horizontal.Mul(0.5)).Sub(c.vertical.Mul(0.5)).Sub(c.orientation.w)
}

func (c *Camera) SetFront(v Vector3) {
	c.orientation.w = v.Unit()
	c.orientation.u = c.orientation.up.Cross(c.orientation.w).Unit()
	c.orientation.v = c.orientation.w.Cross(c.orientation.u)
	c.horizontal = c.orientation.u.Mul(c.viewportWidth)
	c.vertical = c.orientation.v.Mul(c.viewportHeight)
	c.lowerLeftCorner = c.orientation.origin.Sub(c.horizontal.Mul(0.5)).Sub(c.vertical.Mul(0.5)).Sub(c.orientation.w)
}

func (c *Camera) Origin() Vector3 {
	return c.orientation.origin
}

func (c *Camera) Up() Vector3 {
	return c.orientation.up
}

func (c *Camera) W() Vector3 {
	return c.orientation.w
}

// Cast ray in new direction, while keeping origin the same
func (c *Camera) castRayReuse(s, t float64, ray *ray) {
	ray.reuseSameOrigin(c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.orientation.origin))
}

type ray struct {
	origin         Vector3
	direction      Vector3
	dirNormSquared float64

	invDirection Vector3
	sign         [3]int
}

func (r ray) position(t float64) Vector3 {
	magnitude := r.direction.Mul(t)
	return r.origin.Add(magnitude)
}

func newRay(origin Vector3, direction Vector3) ray {
	invDirection := direction.Inverse()
	sign := [3]int{}

	if invDirection.X < 0 {
		sign[0] = 1
	}
	if invDirection.Y < 0 {
		sign[1] = 1
	}
	if invDirection.Z < 0 {
		sign[2] = 1
	}

	dirNormSq := direction.LengthSquared()
	return ray{
		origin:         origin,
		direction:      direction,
		invDirection:   invDirection,
		dirNormSquared: dirNormSq,
		sign:           sign,
	}
}

// Creates a new ray by overriding the already allocated ray
func (r *ray) reuse(origin Vector3, direction Vector3) {
	invDirection := direction.Inverse()
	sign := [3]int{}

	if invDirection.X < 0 {
		sign[0] = 1
	}
	if invDirection.Y < 0 {
		sign[1] = 1
	}
	if invDirection.Z < 0 {
		sign[2] = 1
	}

	dirNormSq := direction.LengthSquared()
	r.origin = origin
	r.direction = direction
	r.invDirection = invDirection
	r.dirNormSquared = dirNormSq
	r.sign = sign
}

// Creates a new ray by overriding the already allocated ray
func (r *ray) reuseSameOrigin(direction Vector3) {
	invDirection := direction.Inverse()
	sign := [3]int{}

	if invDirection.X < 0 {
		sign[0] = 1
	}
	if invDirection.Y < 0 {
		sign[1] = 1
	}
	if invDirection.Z < 0 {
		sign[2] = 1
	}

	dirNormSq := direction.LengthSquared()
	r.direction = direction
	r.invDirection = invDirection
	r.dirNormSquared = dirNormSq
	r.sign = sign
}
