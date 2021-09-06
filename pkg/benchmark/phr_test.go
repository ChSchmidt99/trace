package benchmark

import (
	"github/chschmidt99/pt/demoscenes"
	"github/chschmidt99/pt/pkg/pt"
	"runtime"
	"testing"
)

const (
	AR               = 4.0 / 3.0
	FOV              = 50.0
	RESOLUTION       = 200
	ALPHA            = 0.5
	DELTA            = 6
	BRANCHING_FACTOR = 4
)

func BenchmarkPHRRender(b *testing.B) {
	scene, camera := demoscenes.Bunny(AR, FOV)
	bvh := scene.CompilePHR(ALPHA, DELTA, BRANCHING_FACTOR)
	renderer := pt.NewDefaultRenderer(bvh, camera)
	renderer.Spp = 1
	buff := pt.NewBufferAspect(RESOLUTION, AR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.RenderToBuffer(buff)
	}
}

func BenchmarkPHRBuild(b *testing.B) {
	scene, _ := demoscenes.Bunny(AR, FOV)
	tracables := scene.Tracables()
	builder := pt.NewPHRBuilder(tracables, ALPHA, DELTA, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Build()
	}
}
