package pt

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkPHRBuild(b *testing.B) {
	geometry := ParseFromPath("../../assets/deer.obj")
	tracables := NewMesh(geometry, nil).transformed(IdentityMatrix())
	bounding := Enclosing(tracables)
	aux := LBVH(tracables, bounding, runtime.GOMAXPROCS(0))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildFromAuxilary(tracables, aux, bounding, 0.5, 6, 4, runtime.GOMAXPROCS(0))
	}
}

func TestFindInitialCut(t *testing.T) {
	geometry := ParseFromPath("../../assets/deer.obj")
	tracables := NewMesh(geometry, nil).transformed(IdentityMatrix())
	bounding := Enclosing(tracables)
	aux := LBVH(tracables, bounding, runtime.GOMAXPROCS(0))
	p := phr{
		s:               bounding.surface(),
		alpha:           0.5,
		delta:           6,
		branchingFactor: 4,
		jobs:            make(chan *phrJob),
		primitives:      tracables,
	}
	seq := p.findInitialCut(aux)
	parallel := p.findInitialCutParallel(aux, runtime.GOMAXPROCS(0))
	assert.Equal(t, len(seq.nodes), len(parallel.nodes))
}

/*
func BenchmarkFindInitialCut(b *testing.B) {
	geometry := ParseFromPath("../../assets/local/hairball.obj")
	tracables := NewMesh(geometry, nil).transformed(IdentityMatrix())
	bounding := enclosing(tracables)
	aux := LBVH(tracables, bounding, runtime.GOMAXPROCS(0))
	p := phr{
		s:               bounding.surface(),
		alpha:           0.5,
		delta:           6,
		branchingFactor: 4,
		jobs:            make(chan *phrJob),
		primitives:      tracables,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.findInitialCut(aux)
	}
}

func BenchmarkFindInitialCutParallel(b *testing.B) {
	geometry := ParseFromPath("../../assets/local/hairball.obj")
	tracables := NewMesh(geometry, nil).transformed(IdentityMatrix())
	bounding := enclosing(tracables)
	aux := LBVH(tracables, bounding, runtime.GOMAXPROCS(0))
	p := phr{
		s:               bounding.surface(),
		alpha:           0.5,
		delta:           6,
		branchingFactor: 4,
		jobs:            make(chan *phrJob),
		primitives:      tracables,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.findInitialCutParallel(aux, runtime.GOMAXPROCS(0))
	}
}
*/
