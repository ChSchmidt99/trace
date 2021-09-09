package pt

import (
	"math"
	"testing"
)

func BenchmarkTraversal(b *testing.B) {
	stack := bvhStack{}
	stack.push(&bvhNode{
		bounding: newAABB(NewVector3(0, 0, 0), NewVector3(1, 1, 1)),
		isLeaf:   false,
	})
	ray := newRay(NewVector3(2, 0.5, 0.5), NewVector3(-1, 0, 0))
	tMin := 0.001
	tMax := math.MaxFloat64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node := stack.pop()
		if node == nil {
			panic("Should not be called")
		}
		if node.bounding.intersected(ray, tMin, tMax) {
			if node.isLeaf {
				panic("Should not be called")
			} else {
				stack.push(node)
			}
		}
	}
}

func BenchmarkIntersection(b *testing.B) {
	ray := newRay(NewVector3(2, 2, 2), NewVector3(-1, -1, 1))
	tMin := 0.0001
	tMax := math.MaxFloat64
	tri := NewTriangleWithoutNormals(NewVector3(0, 0, 1), NewVector3(1, 0, 0), NewVector3(0, 1, 0))
	out := hit{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tri.intersected(ray, tMin, tMax, &out)
	}
}

// Single intersection around 7 ns
// Traversal step around 43 ns
