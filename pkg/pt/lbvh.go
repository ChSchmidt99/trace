package pt

import (
	"math"
	"sync"
)

var MORTON_SIZE = uint32(math.Pow(2, 12))

func buildLBVH(prims []Primitive, enclosing aabb, threads int) {
	assignMortonCodes(prims, enclosing, MORTON_SIZE, threads)
}

// TODO: Compare index to pointer approach
type mortonPair struct {
	primIndex  int
	mortonCode uint64
}

func assignMortonCodes(prims []Primitive, enclosing aabb, mortonSize uint32, threads int) []mortonPair {
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

func computeMorton(prim Primitive, morton *Morton, enclosing aabb, mortonSize uint32) uint64 {
	center := prim.bounding().barycenter
	deltaX := math.Abs(enclosing.bounds[0].X - center.X)
	deltaY := math.Abs(enclosing.bounds[0].Y - center.Y)
	deltaZ := math.Abs(enclosing.bounds[0].Z - center.Z)
	xQuantized := uint32(deltaX / (enclosing.width / float64(mortonSize-1)))
	yQuantized := uint32(deltaY / (enclosing.height / float64(mortonSize-1)))
	zQuantized := uint32(deltaZ / (enclosing.depth / float64(mortonSize-1)))
	vec := []uint32{xQuantized, yQuantized, zQuantized}
	// Neglecting error for performance sake
	encoded, _ := morton.Encode(vec)
	return encoded
}
