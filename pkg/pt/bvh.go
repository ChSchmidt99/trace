package pt

import (
	"sync"
	"sync/atomic"
)

type BVH struct {
	root   *bvhNode
	prims  []tracable
	leaves []*bvhNode
}

// Number of intersection tests executed for given ray, including node bounding boxes and leaf primitives
func (bvh *BVH) intersectionTests(ray ray, tMin, tMax float64) int {
	stack := bvhStack{}
	stack.push(bvh.root)
	count := 0
	closest := tMax
	hitOut := hit{}
	for {
		node := stack.pop()
		if node == nil {
			return count
		}
		count++
		if node.bounding.intersected(ray, tMin, closest) {
			if node.isLeaf {
				count += len(node.prims)
				for i := 0; i < len(node.prims); i++ {
					prim := bvh.prims[node.prims[i]]
					if prim.intersected(ray, tMin, closest, &hitOut) {
						closest = hitOut.t
					}
				}
			} else {
				stack.push(node.children...)
			}
		}
	}
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
		go func(pipeline chan *bvhNode, prims []tracable) {
			defer wg.Done()
			for leaf := range pipeline {
				leaf.updateAABB(prims)
			}
		}(jobs, bvh.prims)
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

type bvhNode struct {
	parent       *bvhNode
	prims        []int
	children     []*bvhNode
	bounding     aabb
	m            *sync.Mutex
	childAABBset uint32
	isLeaf       bool
	size         int
}

func newBranch(children int) *bvhNode {
	return &bvhNode{
		children: make([]*bvhNode, children),
		isLeaf:   false,
		m:        &sync.Mutex{},
	}
}

func newLeaf(prims []int) *bvhNode {
	return &bvhNode{
		prims:  prims,
		isLeaf: true,
		size:   1,
	}
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
