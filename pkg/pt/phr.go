package pt

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

// TODO: Refactoring
func DefaultPHR(primitives []tracable) BVH {
	return PHR(primitives, enclosing(primitives), 0.5, 6, 2, runtime.GOMAXPROCS(0))
}

// alpha: How quickly cut size will shrink
// delta: Size of initial cut for d = 0
func PHR(primitives []tracable, enclosing aabb, alpha float64, delta float64, branchingFactor int, threadCount int) BVH {
	auxilaryBVH := LBVH(primitives, enclosing, threadCount)
	return buildFromAuxilary(primitives, auxilaryBVH, enclosing, alpha, delta, branchingFactor, threadCount)
}

func buildFromAuxilary(primitives []tracable, auxilaryBVH BVH, enclosing aabb, alpha float64, delta float64, branchingFactor int, threadCount int) BVH {

	p := phr{
		s:               enclosing.surface(),
		alpha:           alpha,
		delta:           delta,
		branchingFactor: branchingFactor,
		jobs:            make(chan *phrJob, threadCount),
		primitives:      primitives,
	}
	cut := p.findInitialCut(auxilaryBVH)

	wg := sync.WaitGroup{}
	for i := 0; i < threadCount; i++ {
		go func() {
			for job := range p.jobs {
				p.buildSubTree(job, &wg)
			}
		}()
	}

	temp := newBranch(1)
	temp.bounding = enclosing

	wg.Add(1)
	p.queueJob(&phrJob{
		depth:      1,
		cut:        cut,
		parent:     temp,
		childIndex: 0,
	}, &wg)
	wg.Wait()
	close(p.jobs)

	temp.children[0].parent = nil
	return BVH{
		root:  temp.children[0],
		prims: primitives,
	}
}

type phrJob struct {
	depth      int
	cut        *phrCut
	parent     *bvhNode
	childIndex int
}

type phrCut struct {
	nodes    []*bvhNode
	bounding aabb
}

type phr struct {
	s               float64
	alpha, delta    float64
	branchingFactor int
	jobs            chan *phrJob
	primitives      []tracable
}

func (p *phr) buildSubTree(job *phrJob, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(job.cut.nodes) <= 1 {
		job.parent.addChild(job.cut.nodes[0], job.childIndex)
		return
	}

	cuts := make([]*phrCut, 0, p.branchingFactor)
	cuts = append(cuts, job.cut)

	for len(cuts) < p.branchingFactor {
		max := 0
		maxI := 0
		for i, cut := range cuts {
			if len(cut.nodes) > max {
				max = len(cut.nodes)
				maxI = i
			}
		}
		if max <= 1 {
			break
		}
		left, right := p.splitPHRcutAlternative(cuts[maxI])
		cuts[maxI] = p.refined(left, job.depth)
		cuts = append(cuts, p.refined(right, job.depth))
	}

	branch := newBranch(len(cuts))
	branch.parent = job.parent
	branch.bounding = job.cut.bounding
	job.parent.addChild(branch, job.childIndex)

	wg.Add(len(cuts))
	for i, cut := range cuts {
		p.queueJob(&phrJob{
			depth:      job.depth + 1,
			cut:        cut,
			parent:     branch,
			childIndex: i,
		}, wg)
	}
}

func (p *phr) queueJob(job *phrJob, wg *sync.WaitGroup) {
	select {
	case p.jobs <- job:
	default:
		p.buildSubTree(job, wg)
	}
}

func (p *phr) areaThreshold(treeDepth int) float64 {
	return p.s / math.Pow(2, p.alpha*float64(treeDepth)+p.delta)
}

func (p *phr) findInitialCut(lbvh BVH) *phrCut {
	// TODO: Paralellize?
	queue := queue{}
	queue.push(lbvh.root)
	currentCut := &phrCut{
		bounding: lbvh.root.bounding,
	}

	for queue.length > 0 {
		node := queue.popFirst()
		if node.isLeaf {
			currentCut.nodes = append(currentCut.nodes, node)
		} else {
			if node.bounding.surface() > p.areaThreshold(0) {
				for _, child := range node.children {
					queue.push(child)
				}
				continue
			}
			currentCut.nodes = append(currentCut.nodes, node)
		}
	}
	return currentCut
}

func (p *phr) splitPHRcut(cut *phrCut) (left *phrCut, right *phrCut) {

	sort.SliceStable(cut.nodes, func(i, j int) bool {
		return cut.nodes[i].bounding.barycenter.X < cut.nodes[j].bounding.barycenter.X
	})
	minX, iX := minCost(cut.nodes)

	sort.SliceStable(cut.nodes, func(i, j int) bool {
		return cut.nodes[i].bounding.barycenter.Y < cut.nodes[j].bounding.barycenter.Y
	})
	minY, iY := minCost(cut.nodes)

	sort.SliceStable(cut.nodes, func(i, j int) bool {
		return cut.nodes[i].bounding.barycenter.Z < cut.nodes[j].bounding.barycenter.Z
	})
	minZ, iZ := minCost(cut.nodes)

	if minZ < minX && minZ < minY {
		left := cut.nodes[:iZ]
		letfBounding := enclosingSubtrees(left)
		right := cut.nodes[iZ:]
		rightBounding := enclosingSubtrees(right)
		return &phrCut{left, letfBounding}, &phrCut{right, rightBounding}
	}
	if minX < minY && minX < minZ {
		// TODO: Is it more efficient to sort twice, or copy and store sorted slice?
		sort.SliceStable(cut.nodes, func(i, j int) bool {
			return cut.nodes[i].bounding.barycenter.X < cut.nodes[j].bounding.barycenter.X
		})
		left := cut.nodes[:iX]
		letfBounding := enclosingSubtrees(left)
		right := cut.nodes[iX:]
		rightBounding := enclosingSubtrees(right)
		return &phrCut{left, letfBounding}, &phrCut{right, rightBounding}
	} else {
		// TODO: Is it more efficient to sort twice, or copy and store sorted slice? (Same as above)
		sort.SliceStable(cut.nodes, func(i, j int) bool {
			return cut.nodes[i].bounding.barycenter.Y < cut.nodes[j].bounding.barycenter.Y
		})
		left := cut.nodes[:iY]
		letfBounding := enclosingSubtrees(left)
		right := cut.nodes[iY:]
		rightBounding := enclosingSubtrees(right)
		return &phrCut{left, letfBounding}, &phrCut{right, rightBounding}
	}
}

func (p *phr) splitPHRcutAlternative(cut *phrCut) (l *phrCut, r *phrCut) {
	sort.SliceStable(cut.nodes, func(i, j int) bool {
		return cut.nodes[i].bounding.barycenter.X < cut.nodes[j].bounding.barycenter.X
	})
	sorted2 := make([]*bvhNode, len(cut.nodes))
	copy(sorted2, cut.nodes)
	sort.SliceStable(sorted2, func(i, j int) bool {
		return sorted2[i].bounding.barycenter.Y < sorted2[j].bounding.barycenter.Y
	})
	minX, iX := minCost(cut.nodes)
	minY, iY := minCost(sorted2)

	var left []*bvhNode
	var right []*bvhNode
	if minX < minY {
		sort.SliceStable(sorted2, func(i, j int) bool {
			return sorted2[i].bounding.barycenter.Z < sorted2[j].bounding.barycenter.Z
		})
		minZ, iZ := minCost(sorted2)
		if minX < minZ {
			left = cut.nodes[:iX]
			right = cut.nodes[iX:]
		} else {
			left = sorted2[:iZ]
			right = sorted2[iZ:]
		}
	} else {
		sort.SliceStable(cut.nodes, func(i, j int) bool {
			return cut.nodes[i].bounding.barycenter.Z < cut.nodes[j].bounding.barycenter.Z
		})
		minZ, iZ := minCost(cut.nodes)
		if minY < minZ {
			left = sorted2[:iY]
			right = sorted2[iY:]
		} else {
			left = cut.nodes[:iZ]
			right = cut.nodes[iZ:]
		}
	}
	letfBounding := enclosingSubtrees(left)
	rightBounding := enclosingSubtrees(right)
	return &phrCut{left, letfBounding}, &phrCut{right, rightBounding}
}

func (p *phr) refined(cut *phrCut, depth int) *phrCut {
	refinedCut := make([]*bvhNode, 0, len(cut.nodes))
	for _, node := range cut.nodes {
		if node.isLeaf {
			if node.bounding.surface() < p.areaThreshold(depth) {
				refinedCut = append(refinedCut, node)
			} else {
				// TODO: Rethink, is it better to add primitives?
				for _, prim := range node.prims {
					leaf := newLeaf([]int{prim})
					leaf.bounding = p.primitives[prim].bounding()
					refinedCut = append(refinedCut, leaf)
				}
			}
		} else {
			if node.bounding.surface() < p.areaThreshold(depth) {
				// Keep node in cut
				refinedCut = append(refinedCut, node)
			} else {
				// Replace node with children
				refinedCut = append(refinedCut, node.children...)
			}
		}
	}
	return &phrCut{
		nodes:    refinedCut,
		bounding: cut.bounding,
	}
}

func minCost(sortedNodes []*bvhNode) (min float64, splitIndex int) {
	minCost := math.Inf(1)
	minIndex := 0
	for i := 1; i < len(sortedNodes); i++ {
		cost := sahCost(sortedNodes[:i], sortedNodes[i:])
		if cost < minCost {
			minCost = cost
			minIndex = i
		}
	}
	return minCost, minIndex
}

func sahCost(leftCut []*bvhNode, rightCut []*bvhNode) float64 {
	leftSurface := enclosingSubtrees(leftCut).surface()
	leftNodeCount := nodeCount(leftCut)
	rightSurface := enclosingSubtrees(rightCut).surface()
	rightNodeCount := nodeCount(rightCut)
	return leftSurface*float64(leftNodeCount) + rightSurface*float64(rightNodeCount)
}

// TODO: Make size better available?
func nodeCount(subtrees []*bvhNode) int {
	sum := 0
	for _, node := range subtrees {
		sum += node.subtreeSize()
	}
	return sum
}

// TODO: There is probably a more efficient way using channels
type queue struct {
	length int
	first  *queueEntry
}

func (q *queue) push(node *bvhNode) {
	entry := queueEntry{value: node}
	if q.first == nil {
		q.first = &entry
		q.length = 1
		return
	}
	q.length++
	q.first.append(&entry)
}

func (q *queue) popFirst() *bvhNode {
	if q.first == nil {
		q.length = 0
		return nil
	}
	q.length--
	entry := q.first
	if entry.next != nil {
		q.first = entry.next
	} else {
		q.first = nil
	}
	return entry.value
}

type queueEntry struct {
	next  *queueEntry
	value *bvhNode
}

func (e *queueEntry) append(entry *queueEntry) {
	if e.next == nil {
		e.next = entry
		return
	}
	e.next.append(entry)
}
