package pt

import (
	"runtime"
	"sync"
)

// Temporary BVH placeholder for testing

type BVH struct {
	root  *bvhNode
	prims []Primitive
}

func NewBVH(prims []Primitive) BVH {
	/*
		indeces := make([]int, len(prims))
		for i := 0; i < len(prims); i++ {
			indeces[i] = i
		}
		return BVH{
			root: &bvhNode{
				prims:    indeces,
				bounding: enclosing(prims),
			},
			prims: prims,
		}
	*/
	return buildLBVH(prims, enclosing(prims), runtime.GOMAXPROCS(0))
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
			if node.isLeaf {
				for i := 0; i < len(node.prims); i++ {
					prim := bvh.prims[node.prims[i]]
					if prim.intersected(ray, tMin, closest, hitOut) {
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

// TODO: Rework bounding box update
func (bvh *BVH) updateBounding(threads int) {
	leaves := bvh.collectLeaves()
	wg := sync.WaitGroup{}
	wg.Add(threads)
	jobs := make(chan *bvhNode)
	for i := 0; i < threads; i++ {
		go func(pipeline chan *bvhNode, prims []Primitive) {
			defer wg.Done()
			for leaf := range pipeline {
				leaf.updateAABB(prims)
			}
		}(jobs, bvh.prims)
	}
	for _, leaf := range leaves {
		jobs <- leaf
	}
	close(jobs)
	wg.Wait()
}

// TODO: Can this be parallelized or replaced?
func (bvh *BVH) collectLeaves() []*bvhNode {
	acc := make([]*bvhNode, 0)
	bvh.root.collectLeaves(&acc)
	return acc
}

// TODO: Test interface bvh and compare performance
// TODO: Test performance of array instead of node pointers
type bvhNode struct {
	parent       *bvhNode
	prims        []int
	children     []*bvhNode
	bounding     aabb
	m            *sync.Mutex
	childAABBset int
	isLeaf       bool
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
	}
}

func (node *bvhNode) addChild(child *bvhNode, index int) {
	node.children[index] = child
	child.parent = node
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

// TODO: rethink function, maybe using a channel and pushing all ready nodes makes more sense
func (node *bvhNode) updateAABB(primitives []Primitive) {
	if node.isLeaf {
		node.bounding = enclosingSlice(node.prims, primitives)
		// TODO: use atomic.CompareAndSwapUint32
		node.parent.m.Lock()
		node.parent.childAABBset++
		if node.parent.childAABBset%len(node.parent.children) == 0 {
			node.parent.m.Unlock()
			node.parent.updateAABB(primitives)
			return
		}
		node.parent.m.Unlock()
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
	node.parent.m.Lock()
	node.parent.childAABBset++
	if node.parent.childAABBset%2 == 0 {
		node.parent.m.Unlock()
		node.parent.updateAABB(primitives)
		return
	}
	node.parent.m.Unlock()
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
