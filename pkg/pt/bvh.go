package pt

// Temporary BVH placeholder for testing

type BVH struct {
	root *bvhNode
}

func NewBVH(prims []Primitive) BVH {
	return BVH{
		&bvhNode{
			prims:    prims,
			bounding: enclosing(prims),
		},
	}
}

func (bvh *BVH) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	stack := bvhStack{}
	stack.push(bvh.root)
	didHit := false
	closest := tMax
	for {
		node := stack.pop()
		if node == nil {
			return didHit
		}
		if node.bounding.intersected(ray, tMin, closest) {
			if node.isLeaf() {
				for i := 0; i < len(node.prims); i++ {
					if node.prims[i].intersected(ray, tMin, closest, hitOut) {
						didHit = true
						closest = hitOut.t
					}
				}
			} else {
				stack.push(node.children...)
			}
		}
	}
}

// TODO: Test interface bvh and compare performance
// TODO: Add constructor + children capacity according to branching factor
type bvhNode struct {
	prims    []Primitive
	children []*bvhNode
	bounding aabb
}

// TODO: Check performance of function
func (node *bvhNode) isLeaf() bool {
	return len(node.children) == 0
}

// TODO: Reuse stack nodes?
type bvhStack struct {
	top *bvhStackNode
}

type bvhStackNode struct {
	val  *bvhNode
	prev *bvhStackNode
}

// Fist element of the slice will be pushed to stack first
func (s *bvhStack) push(vals ...*bvhNode) {
	for _, val := range vals {
		s.top = &bvhStackNode{
			val:  val,
			prev: s.top,
		}
	}
}

func (s *bvhStack) pop() *bvhNode {
	if s.top == nil {
		return nil
	}
	n := s.top
	s.top = n.prev
	return n.val
}
