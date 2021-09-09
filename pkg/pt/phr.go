package pt

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

type PhrBuilder struct {
	Alpha           float64 // How quickly cut size will shrink
	Delta           int     // Determines size of initial cut
	BranchingFactor int
	Threshold       AreaThreshold
	Split           SplitFunction
	jobs            chan *phrJob
	threadCount     int
	primitives      []tracable
	surface         float64
}

func NewDefaultBuilder(primitives []tracable) PhrBuilder {
	return NewPHRBuilder(primitives, 0.5, 6, 4, runtime.GOMAXPROCS(0))
}

func NewPHRBuilder(primitives []tracable, alpha float64, delta int, branchingFactor int, threadCount int) PhrBuilder {
	return PhrBuilder{
		Alpha:           alpha,
		Delta:           delta,
		BranchingFactor: branchingFactor,
		Threshold:       DefaultThreshold,
		Split:           SweepSAH,
		primitives:      primitives,
		threadCount:     threadCount,
	}
}

func (p PhrBuilder) Build() BVH {
	auxilaryBVH := LBVH(p.primitives, enclosing(p.primitives), p.threadCount)
	return p.BuildFromAuxilary(auxilaryBVH)
}

func (p PhrBuilder) BuildFromAuxilary(auxilaryBVH BVH) BVH {
	p.surface = auxilaryBVH.root.bounding.surface()
	// Determine initial cut
	cut := p.findInitialCut(auxilaryBVH, p.threadCount)
	// Start workers
	wg := sync.WaitGroup{}
	p.jobs = make(chan *phrJob, p.threadCount)
	for i := 0; i < p.threadCount; i++ {
		go func() {
			for job := range p.jobs {
				p.buildSubTree(job, &wg)
			}
		}()
	}

	// Temporary branch as a starting point, will be discared afterwards
	temp := newBranch(1)
	temp.bounding = auxilaryBVH.root.bounding
	wg.Add(1)

	// Start initial job
	p.queueJob(&phrJob{
		depth:      1,
		cut:        cut,
		parent:     temp,
		childIndex: 0,
	}, &wg)

	// Wait until tree is built
	wg.Wait()
	close(p.jobs)

	temp.children[0].parent = nil
	return BVH{
		root:  temp.children[0],
		prims: p.primitives,
	}
}

type phrJob struct {
	depth      int
	cut        phrCut
	parent     *bvhNode
	childIndex int
}

func (p PhrBuilder) queueJob(job *phrJob, wg *sync.WaitGroup) {
	// Pass job to channel if there is capacity left, otherwise process it directly
	select {
	case p.jobs <- job:
	default:
		p.buildSubTree(job, wg)
	}
}

func (p PhrBuilder) buildSubTree(job *phrJob, wg *sync.WaitGroup) {
	defer wg.Done()

	// Termination criteria: if only one node is left in the cut, collect its primitives and cobine into new leaf
	/*
		if len(job.cut.nodes) <= 1 {
			leaves := make([]*bvhNode, 0)
			job.cut.nodes[0].collectLeaves(&leaves)
			prims := make([]int, 0, len(leaves))
			for _, leaf := range leaves {
				prims = append(prims, leaf.prims...)
			}
			leaf := newLeaf(prims)
			leaf.bounding = job.cut.bounding
			job.parent.addChild(leaf, job.childIndex)
			return
		}
	*/
	if len(job.cut.nodes) <= 1 {
		// Adding the child this way means, that parent pointers will be wrong
		// This can't be avoided, as the auxiliary BVH would be destroyed otherwise
		// A solution would be copying the subtree, however, this is not necessarry,
		// as parent pointers are not used in the final BVH
		job.parent.children[job.childIndex] = job.cut.nodes[0]
		return
	}

	cuts := make([]phrCut, 1, p.BranchingFactor)
	cuts[0] = job.cut

	// Keep splitting cut until enough nodes to branch the tree are found
	for len(cuts) < p.BranchingFactor {
		// Find the biggest cut
		max := 0
		maxI := 0
		for i, cut := range cuts {
			if len(cut.nodes) > max {
				max = len(cut.nodes)
				maxI = i
			}
		}
		// If the biggest cut has size = 1, no more cuts can be split => break
		if max <= 1 {
			break
		}
		// Split biggest cut
		left, right := p.Split(cuts[maxI])
		cuts[maxI] = p.refined(left, job.depth)
		cuts = append(cuts, p.refined(right, job.depth))
	}

	// Create a new BVH branch
	branch := newBranch(len(cuts))
	branch.parent = job.parent
	branch.bounding = job.cut.bounding
	job.parent.addChild(branch, job.childIndex)

	// Queue all new children to be processed by this or any other thread
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

func (p PhrBuilder) findInitialCut(auxilary BVH, threadCount int) phrCut {
	queue := make(chan *bvhNode, 1024)
	cut := phrCut{
		bounding: auxilary.root.bounding,
	}
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < threadCount; i++ {
		go func() {
			for node := range queue {
				if node.isLeaf {
					// Add node to cut, if it is a leaf
					m.Lock()
					cut.nodes = append(cut.nodes, node)
					m.Unlock()
					wg.Done()
				} else {
					// Add children to queue, if
					if node.bounding.surface() > p.Threshold(p.surface, p.Alpha, p.Delta, 0) {
						wg.Add(len(node.children) - 1)
						for _, child := range node.children {
							queue <- child
						}
						continue
					}
					m.Lock()
					cut.nodes = append(cut.nodes, node)
					m.Unlock()
					wg.Done()
				}
			}
		}()
	}
	queue <- auxilary.root
	wg.Add(1)
	wg.Wait()
	close(queue)
	return cut
}

func (p PhrBuilder) refined(cut phrCut, depth int) phrCut {
	refinedCut := make([]*bvhNode, 0, len(cut.nodes))
	for _, node := range cut.nodes {
		if node.isLeaf {
			if node.bounding.surface() < p.Threshold(p.surface, p.Alpha, p.Delta, depth) {
				refinedCut = append(refinedCut, node)
			} else {
				for _, prim := range node.prims {
					leaf := newLeaf([]int{prim})
					leaf.bounding = p.primitives[prim].bounding()
					refinedCut = append(refinedCut, leaf)
				}
			}
		} else {
			if node.bounding.surface() < p.Threshold(p.surface, p.Alpha, p.Delta, depth) {
				// Keep node in cut
				refinedCut = append(refinedCut, node)
			} else {
				// Replace node with children
				refinedCut = append(refinedCut, node.children...)
			}
		}
	}
	return phrCut{
		nodes:    refinedCut,
		bounding: cut.bounding,
	}
}

type AreaThreshold func(surface float64, alpha float64, delta int, depth int) float64

func DefaultThreshold(surface float64, alpha float64, delta int, depth int) float64 {
	return surface / math.Pow(2, alpha*float64(depth)+float64(delta))
}

type phrCut struct {
	nodes    []*bvhNode
	bounding aabb
}

type SplitFunction func(phrCut) (phrCut, phrCut)

func SweepSAH(cut phrCut) (l phrCut, r phrCut) {
	if len(cut.nodes) == 2 {
		return phrCut{
				nodes:    cut.nodes[:1],
				bounding: cut.nodes[0].bounding,
			}, phrCut{
				nodes:    cut.nodes[1:],
				bounding: cut.nodes[1].bounding,
			}
	}

	// Sort along x and y axis using two separate slices
	sort.SliceStable(cut.nodes, func(i, j int) bool {
		return cut.nodes[i].bounding.barycenter.X < cut.nodes[j].bounding.barycenter.X
	})
	sorted2 := make([]*bvhNode, len(cut.nodes))
	copy(sorted2, cut.nodes)
	sort.SliceStable(sorted2, func(i, j int) bool {
		return sorted2[i].bounding.barycenter.Y < sorted2[j].bounding.barycenter.Y
	})

	xSAH := minCost(cut.nodes)
	ySAH := minCost(sorted2)
	// Keep the sorted slice with lower cost and override the other by sorting along z axis
	// Finally, return the split with the lowest cost

	if xSAH.cost < ySAH.cost {
		sort.SliceStable(sorted2, func(i, j int) bool {
			return sorted2[i].bounding.barycenter.Z < sorted2[j].bounding.barycenter.Z
		})
		zSAH := minCost(sorted2)
		if xSAH.cost < zSAH.cost {
			return xSAH.left, xSAH.right
		} else {
			return zSAH.left, zSAH.right
		}
	} else {
		sort.SliceStable(cut.nodes, func(i, j int) bool {
			return cut.nodes[i].bounding.barycenter.Z < cut.nodes[j].bounding.barycenter.Z
		})
		zSAH := minCost(cut.nodes)
		if ySAH.cost < zSAH.cost {
			return ySAH.left, ySAH.right
		} else {
			return zSAH.left, zSAH.right
		}
	}
}

type sah struct {
	left  phrCut
	right phrCut
	cost  float64
}

// Uses SAH to compute the best split
func minCost(sortedNodes []*bvhNode) sah {
	min := sah{
		cost: math.Inf(1),
	}
	// Compute and track right costs by incrementally extending bounding box
	SaRight := sortedNodes[len(sortedNodes)-1].bounding
	rightCosts := make([]float64, len(sortedNodes))
	rightCuts := make([]phrCut, len(sortedNodes))
	nodeCount := 0
	for i := len(sortedNodes) - 1; i > 0; i-- {
		SaRight = SaRight.add(sortedNodes[i].bounding)
		nodeCount += sortedNodes[i].subtreeSize()
		rightCosts[i] = SaRight.surface() * float64(nodeCount)
		rightCuts[i] = phrCut{
			nodes:    sortedNodes[i:],
			bounding: SaRight,
		}
	}

	// Incrementally extend left box and use tracked right costs to compute full SAH cost
	nodeCount = sortedNodes[0].subtreeSize()
	SaLeft := sortedNodes[0].bounding
	for i := 1; i < len(sortedNodes); i++ {
		cost := rightCosts[i] + SaLeft.surface()*float64(nodeCount)
		if cost < min.cost {
			min.cost = cost
			min.left = phrCut{
				nodes:    sortedNodes[:i],
				bounding: SaLeft,
			}
			min.right = rightCuts[i]
		}
		SaLeft = SaLeft.add(sortedNodes[i].bounding)
		nodeCount += sortedNodes[i].subtreeSize()
	}
	return min
}
