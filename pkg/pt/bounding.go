package pt

// TODO: Compare to pointer aabb
type aabb struct {
	bounds [2]Vector3 // bounds[0] = min, bounds[1] = max
}

func newAABB(min, max Vector3) aabb {
	return aabb{
		bounds: [2]Vector3{min, max},
	}
}

func enclosing(primitives []Primitive) aabb {
	enclosing := primitives[0].bounding()
	for i := 1; i < len(primitives); i++ {
		enclosing = enclosing.add(primitives[i].bounding())
	}
	return enclosing
}

func (a aabb) add(b aabb) aabb {
	return newAABB(MinVec(a.bounds[0], b.bounds[0]), MaxVec(a.bounds[1], b.bounds[1]))
}

func (aabb aabb) intersected(ray ray, tMin, tMax float64) bool {

	tXmin := (aabb.bounds[ray.sign[0]].X - ray.origin.X) * ray.invDirection.X
	tXmax := (aabb.bounds[1-ray.sign[0]].X - ray.origin.X) * ray.invDirection.X
	tYmin := (aabb.bounds[ray.sign[1]].Y - ray.origin.Y) * ray.invDirection.Y
	tYmax := (aabb.bounds[1-ray.sign[1]].Y - ray.origin.Y) * ray.invDirection.Y

	if tXmin > tYmax || tYmin > tXmax {
		return false
	}
	if tYmin > tXmin {
		tXmin = tYmin
	}
	if tYmax < tXmax {
		tXmax = tYmax
	}

	tZmin := (aabb.bounds[ray.sign[2]].Z - ray.origin.Z) * ray.invDirection.Z
	tZmax := (aabb.bounds[1-ray.sign[2]].Z - ray.origin.Z) * ray.invDirection.Z

	if tXmin > tZmax || tZmin > tXmax {
		return false
	}

	// Check if the ray origin is in front of the box or behind
	if tZmin > tXmin {
		tXmin = tZmin
	}

	if tZmax < tXmax {
		tXmax = tZmax
	}

	if tXmax < tMin || tXmin > tMax {
		return false
	}

	return true
}
