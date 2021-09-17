package pt

type aabb struct {
	bounds     [2]Vector3 // bounds[0] = min, bounds[1] = max
	width      float64
	height     float64
	depth      float64
	barycenter Vector3
}

func newAABB(min, max Vector3) aabb {
	bouding := aabb{
		bounds: [2]Vector3{min, max},
	}
	bouding.update()
	return bouding
}

func enclosing(primitives []tracable) aabb {
	enclosing := primitives[0].bounding()
	for i := 1; i < len(primitives); i++ {
		enclosing = enclosing.add(primitives[i].bounding())
	}
	return enclosing
}

func enclosingSlice(indeces []int, primitives []tracable) aabb {
	enclosing := primitives[indeces[0]].bounding()
	for i := 1; i < len(indeces); i++ {
		prim := primitives[indeces[i]]
		enclosing = enclosing.add(prim.bounding())
	}
	return enclosing
}

func (a *aabb) update() {
	min := a.bounds[0]
	max := a.bounds[1]
	a.barycenter = min.Add(max).Mul(1.0 / 2.0)
	a.width = max.X - min.X
	a.height = max.Y - min.Y
	a.depth = max.Z - min.Z
}

func (a *aabb) surface() float64 {
	return 2*a.width*a.height + 2*a.width*a.depth + 2*a.height*a.depth
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
