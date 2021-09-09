package pt

import (
	"sync"
	"sync/atomic"
)

// TODO: Set constants
const (
	INTERSECTION_COST = 7  // roughly approximated cost of intersection calculation
	TRAVERSAL_COST    = 43 // cost of traversal relative to intersection cost
)

type BVH struct {
	root   *bvhNode
	prims  []tracable
	leaves []*bvhNode
}

// Number of intersection tests executed for given ray, including node bounding boxes and leaf primitives
func (bvh *BVH) traversalSteps(ray ray, tMin, tMax float64) int {
	stack := bvhStack{}
	stack.push(bvh.root)
	count := 0
	closest := tMax
	for {
		node := stack.pop()
		if node == nil {
			return count
		}
		if node.bounding.intersected(ray, tMin, closest) {
			if !node.isLeaf {
				stack.push(node.children...)
				count++
			}
		}
	}
}

func (bvh *BVH) cost() float64 {
	return bvh.root.cost()
}

func (bvh *BVH) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	stack := bvhStack{}
	stack.push(bvh.root)
	didHit := false
	hitOut.t = tMax
	for {
		node := stack.pop()
		if node == nil {
			return didHit
		}
		if node.bounding.intersected(ray, tMin, hitOut.t) {
			if node.isLeaf {
				for i := 0; i < len(node.prims); i++ {
					prim := bvh.prims[node.prims[i]]
					if prim.intersected(ray, tMin, hitOut.t, hitOut) {
						didHit = true
					}
				}
			} else {
				stack.push(node.children...)
			}
		}
	}
}

func (bvh *BVH) updateBounding(threads int) {
	wg := sync.WaitGroup{}
	wg.Add(threads)
	jobs := make(chan *bvhNode)
	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			for leaf := range jobs {
				leaf.updateAABB(bvh.prims)
			}
		}()
	}
	for _, leaf := range bvh.leaves {
		jobs <- leaf
	}
	close(jobs)
	wg.Wait()
}

func (bvh *BVH) storeLeaves() {
	bvh.leaves = make([]*bvhNode, 0)
	bvh.root.collectLeaves(&bvh.leaves)
}

func (bvh *BVH) size() int {
	return bvh.root.subtreeSize()
}

type bvhNode struct {
	parent       *bvhNode
	prims        []int
	children     []*bvhNode
	bounding     aabb
	childAABBset uint32 // used as atomic counter when updating AABB
	isLeaf       bool
	size         int
}

func (node *bvhNode) initLeaf(prims []int) {
	node.isLeaf = true
	node.size = 1
	node.prims = prims
}

func (node *bvhNode) initBranch(children int) {
	node.isLeaf = false
	node.children = make([]*bvhNode, children)
}

func (node *bvhNode) addChild(child *bvhNode, index int) {
	node.children[index] = child
	child.parent = node
}

func (node *bvhNode) subtreeSize() int {
	if node.size == 0 {
		node.size = 1
		for _, child := range node.children {
			node.size += child.subtreeSize()
		}
	}
	return node.size
}

func (node *bvhNode) collectLeaves(acc *[]*bvhNode) {
	if node.isLeaf {
		*acc = append(*acc, node)
		return
	}
	for _, child := range node.children {
		child.collectLeaves(acc)
	}
}

func (node *bvhNode) updateAABB(primitives []tracable) {
	if node.isLeaf {
		node.bounding = enclosingSlice(node.prims, primitives)
		// Atomic counter. after all child bounding boxes have been computed the parents bounding box can be calculated
		if atomic.AddUint32(&node.parent.childAABBset, 1)%uint32(len(node.parent.children)) == 0 {
			node.parent.updateAABB(primitives)
		}
		return
	}
	node.bounding = node.children[0].bounding
	for i := 1; i < len(node.children); i++ {
		node.bounding = node.bounding.add(node.children[i].bounding)
	}
	node.bounding.update()
	if node.parent == nil {
		return
	}
	if atomic.AddUint32(&node.parent.childAABBset, 1)%uint32(len(node.parent.children)) == 0 {
		node.parent.updateAABB(primitives)
	}
}

// Measurement of the bvh quality
func (node *bvhNode) cost() float64 {
	if node.isLeaf {
		return INTERSECTION_COST * float64(len(node.prims))
	} else {
		childCosts := 0.0
		for _, child := range node.children {
			probability := child.bounding.surface() / node.bounding.surface()
			childCosts += probability * child.cost()
		}
		return TRAVERSAL_COST + childCosts
	}
}

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
