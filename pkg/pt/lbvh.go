package pt

import (
	"math"
	"runtime"
	"sync"
)

var MORTON_SIZE = uint32(math.Pow(2, 20))

const BUCKET_COUNT = 4096

func DefaultLBVH(prims []tracable) BVH {
	return LBVH(prims, enclosing(prims), runtime.GOMAXPROCS(0))
}

func LBVH(prims []tracable, enclosing aabb, threads int) BVH {
	pairs := assignMortonCodes(prims, enclosing, MORTON_SIZE, threads)
	maxMorton := uint64(math.Pow(float64(MORTON_SIZE), 3))
	pairs = sortMortonPairs(pairs, BUCKET_COUNT, maxMorton, threads)
	bvh := constructLBVH(pairs, MORTON_SIZE, threads)
	bvh.prims = prims
	bvh.storeLeaves()
	bvh.updateBounding(threads)
	return bvh
}

type mortonPair struct {
	primIndex  int
	mortonCode uint64
}

// Iterates over all primitives in parallel and assigns morton codes
// TODO: Check if returning a corresponding morton code slice is faster than using structs!
func assignMortonCodes(prims []tracable, enclosing aabb, mortonSize uint32, threads int) []mortonPair {
	pairs := make([]mortonPair, len(prims))
	batchSize := int(math.Ceil(float64(len(prims)) / float64(threads)))
	wg := sync.WaitGroup{}
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		start := i * batchSize
		end := start + batchSize
		if start >= len(prims) {
			wg.Done()
			continue
		}
		if end > len(prims) {
			end = len(prims)
		}
		go func() {
			morton := NewMorton(3, mortonSize)
			for j := start; j < end; j++ {
				code := computeMorton(prims[j], morton, enclosing, mortonSize)
				pairs[j] = mortonPair{
					primIndex:  j,
					mortonCode: code,
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return pairs
}

// Computes a morton code according to the quantized primitive centroid
func computeMorton(prim tracable, morton *Morton, enclosing aabb, mortonSize uint32) uint64 {
	center := prim.bounding().barycenter
	deltaX := math.Abs(enclosing.bounds[0].X - center.X)
	deltaY := math.Abs(enclosing.bounds[0].Y - center.Y)
	deltaZ := math.Abs(enclosing.bounds[0].Z - center.Z)
	xQuantized := uint32(deltaX / (enclosing.width / float64(mortonSize-1)))
	yQuantized := uint32(deltaY / (enclosing.height / float64(mortonSize-1)))
	zQuantized := uint32(deltaZ / (enclosing.depth / float64(mortonSize-1)))
	vec := []uint32{xQuantized, yQuantized, zQuantized}
	encoded, _ := morton.Encode(vec)
	return encoded
}

// Constructs BVH by inserting sorted morton primitive pairs into a binary radix tree
func constructLBVH(pairs []mortonPair, mortonSize uint32, threads int) BVH {
	splitMask := uint64(math.Pow(float64(mortonSize), 3)) / 2
	wg := sync.WaitGroup{}
	wg.Add(len(pairs))
	queue := lbvhWorkerQueue{
		jobs: make(chan *lbvhJob, threads),
		wg:   &wg,
	}

	// Start workers, each worker will find a split in its given interval and spawn 2 new jobs
	for i := 0; i < threads; i++ {
		go func(q *lbvhWorkerQueue) {
			for job := range queue.jobs {
				job.process(q)
			}
		}(&queue)
	}

	initialNode := newBranch(2)

	initialJob := lbvhJob{
		pairs:      pairs,
		splitMask:  splitMask,
		parent:     initialNode,
		childIndex: 0,
	}
	queue.add(&initialJob)
	wg.Wait()

	root := initialNode.children[0]
	root.parent = nil
	return BVH{
		root: root,
	}
}

type lbvhJob struct {
	pairs      []mortonPair
	splitMask  uint64
	parent     *bvhNode
	childIndex int
}

type lbvhWorkerQueue struct {
	jobs chan *lbvhJob
	wg   *sync.WaitGroup
}

func (queue *lbvhWorkerQueue) add(job *lbvhJob) {
	select {
	case queue.jobs <- job:
	default:
		job.process(queue)
	}
}

func (job *lbvhJob) process(queue *lbvhWorkerQueue) {
	if isLeaf(job.pairs) {
		indeces := make([]int, len(job.pairs))
		queue.wg.Add(1)
		for i, pair := range job.pairs {
			indeces[i] = pair.primIndex
			queue.wg.Done()
		}
		job.parent.addChild(newLeaf(indeces), job.childIndex)
		queue.wg.Done()
		return
	}

	// Find the split in the given interval where the most significant bit first changes
	splitIndex := findSplit(job.pairs, job.splitMask)

	// TODO: Remove and rollup singletons after?
	// If there is no split, only spawn one job, which makes pruning step afterwards obsolete
	if splitIndex == 0 || splitIndex == len(job.pairs) {
		queue.add(&lbvhJob{
			pairs:      job.pairs,
			splitMask:  job.splitMask >> 1,
			parent:     job.parent,
			childIndex: job.childIndex,
		})
		return
	}
	// Create a new branch and spawn new jobs for both children
	branch := newBranch(2)
	job.parent.addChild(branch, job.childIndex)
	left := job.pairs[:splitIndex]
	right := job.pairs[splitIndex:]
	queue.add(&lbvhJob{
		pairs:      left,
		splitMask:  job.splitMask >> 1,
		parent:     branch,
		childIndex: 0,
	})
	queue.add(&lbvhJob{
		pairs:      right,
		splitMask:  job.splitMask >> 1,
		parent:     branch,
		childIndex: 1,
	})
}

func isLeaf(pairs []mortonPair) bool {
	return pairs[0].mortonCode == pairs[len(pairs)-1].mortonCode
}

// Binary search to find index of primMortonPair that first exceeds splitMortonCode
func findSplit(pairs []mortonPair, splitMask uint64) int {
	if (pairs[0].mortonCode & splitMask) > 0 {
		return 0
	}
	l := 0
	r := len(pairs) - 1
	for l <= r {
		m := (l + r) / 2
		if (pairs[m].mortonCode & splitMask) == 0 {
			// Continue search in right half
			l = m + 1
		} else if (pairs[m].mortonCode&splitMask) > 0 && (pairs[m-1].mortonCode&splitMask) > 0 {
			// Continue search in left half
			r = m - 1
		} else {
			return m
		}
	}
	return len(pairs)
}
