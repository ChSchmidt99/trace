package pt

// Temporary BVH placeholder for testing

type BVH struct {
	root *bvhNode
}

func NewBVH(prims []Primitive) BVH {
	return BVH{
		&bvhNode{
			prims: prims,
		},
	}
}

func (bvh *BVH) intersected(ray *ray, tMin, tMax float64) *hit {
	return bvh.root.intersected(ray, tMin, tMax)
}

type bvhNode struct {
	prims []Primitive
}

func (node *bvhNode) intersected(ray *ray, tMin, tMax float64) *hit {
	closest := tMax
	var record *hit
	for _, prim := range node.prims {
		if hit := prim.intersected(ray, tMin, closest); hit != nil {
			closest = hit.t
			record = hit
		}
	}
	return record
}
