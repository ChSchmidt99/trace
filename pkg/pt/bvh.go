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

func (bvh *BVH) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	return bvh.root.intersected(ray, tMin, tMax, hitOut)
}

type bvhNode struct {
	prims []Primitive
}

func (node *bvhNode) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	closest := tMax
	didHit := false
	for i := 0; i < len(node.prims); i++ {
		if node.prims[i].intersected(ray, tMin, closest, hitOut) {
			didHit = true
			closest = hitOut.t
		}
	}
	return didHit
}
