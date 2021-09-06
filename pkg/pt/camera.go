package pt

import (
	"math"
)

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

// Cast ray in new direction, while keeping origin the same
func (c *Camera) castRayReuse(s, t float64, ray *ray) {
	ray.reuseSameOrigin(c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.orientation.origin))
}

type ray struct {
	origin    Vector3
	direction Vector3
	//unitDir        Vector3 // TODO: unit dir used?
	dirNormSquared float64

	invDirection Vector3
	sign         [3]int
}

func (r ray) Position(t float64) Vector3 {
	magnitude := r.direction.Mul(t)
	return r.origin.Add(magnitude)
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
	//r.unitDir = direction.Unit()
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
	//r.unitDir = direction.Unit()
	r.dirNormSquared = dirNormSq
	r.sign = sign
}
