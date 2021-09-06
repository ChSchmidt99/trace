package benchmark

import (
	"github/chschmidt99/pt/demoscenes"
	"github/chschmidt99/pt/pkg/pt"
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
	buff := pt.NewBufferAspect(RESOLUTION, AR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.RenderToBuffer(buff)
	}
}
