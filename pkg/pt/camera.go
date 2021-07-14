package pt

import (
	"math"
)

type Ray struct {
	origin         Vector3
	direction      Vector3
	unitDir        Vector3 // TODO: only use unitDir?
	dirNormSquared float64

	invDirection Vector3
	sign         [3]int
}

// TODO: Can old rays be reused instead of creating new ones?
func NewRay(origin Vector3, direction Vector3) *Ray {
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

	dirNorm := direction.Length()

	return &Ray{
		origin:         origin,
		direction:      direction,
		invDirection:   invDirection,
		unitDir:        direction.Unit(),
		dirNormSquared: dirNorm * dirNorm,
		sign:           sign,
	}
}

func (r *Ray) Position(t float64) Vector3 {
	magnitude := r.direction.Mul(t)
	return r.origin.Add(magnitude)
}

type CameraTransformation struct {
	LookFrom Vector3
	LookAt   Vector3
	Up       Vector3
}

type CameraOrientation struct {
	origin Vector3
	up     Vector3
	w      Vector3
	u      Vector3
	v      Vector3
}

func NewOrientation(lookFrom Vector3, lookAt Vector3, up Vector3) CameraOrientation {
	w := lookFrom.Sub(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)
	return CameraOrientation{
		origin: lookFrom,
		w:      lookFrom.Sub(lookAt).Unit(),
		up:     up,
		u:      u,
		v:      v,
	}
}

type Camera struct {
	orientation CameraOrientation

	viewportWidth  float64
	viewportHeight float64

	lowerLeftCorner Vector3
	horizontal      Vector3
	vertical        Vector3
}

func NewCamera(aspectRatio float64, fov float64, transform CameraTransformation) *Camera {
	cam := new(Camera)
	cam.orientation = NewOrientation(transform.LookFrom, transform.LookAt, transform.Up)
	cam.viewportHeight = 2.0 * math.Tan(DegreesToRadians(fov)/2)
	cam.viewportWidth = aspectRatio * cam.viewportHeight
	cam.horizontal = cam.orientation.u.Mul(cam.viewportWidth)
	cam.vertical = cam.orientation.v.Mul(cam.viewportHeight)
	cam.lowerLeftCorner = cam.orientation.origin.Sub(cam.horizontal.Mul(0.5)).Sub(cam.vertical.Mul(0.5)).Sub(cam.orientation.w)
	return cam
}

// TODO: Split up camera movement and castRay
func (cam *Camera) translate(v Vector3) {
	cam.orientation.origin = cam.orientation.origin.Add(v)
	cam.lowerLeftCorner = cam.orientation.origin.Sub(cam.horizontal.Mul(0.5)).Sub(cam.vertical.Mul(0.5)).Sub(cam.orientation.w)
}

func (cam *Camera) setFront(v Vector3) {
	cam.orientation.w = v.Unit()
	cam.orientation.u = cam.orientation.up.Cross(cam.orientation.w).Unit()
	cam.orientation.v = cam.orientation.w.Cross(cam.orientation.u)
	cam.horizontal = cam.orientation.u.Mul(cam.viewportWidth)
	cam.vertical = cam.orientation.v.Mul(cam.viewportHeight)
	cam.lowerLeftCorner = cam.orientation.origin.Sub(cam.horizontal.Mul(0.5)).Sub(cam.vertical.Mul(0.5)).Sub(cam.orientation.w)
}

func (c *Camera) castRay(s, t float64) *Ray {
	return NewRay(c.orientation.origin, c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.orientation.origin))
}
