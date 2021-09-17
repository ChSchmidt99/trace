package benchmark

import (
	"github/chschmidt99/pt/pkg/demoscenes"
	"github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
	"runtime"
	"strconv"
	"testing"
)

const (
	AR         = 1.0
	FOV        = 50.0
	FRAME_SIZE = 256
	//FRAME_SIZE       = 512
	PHR_FAST_ALPHA   = 0.5
	PHR_FAST_DELTA   = 6.0
	PHR_HQ_ALPHA     = 0.55
	PHR_HQ_DELTA     = 9.0
	BRANCHING_FACTOR = 2
)

var demoScene = demoscenes.Bunny(AR, FOV)

//var demoScene = demoscenes.SanMiguel(AR, FOV)

var alpha = PHR_FAST_ALPHA
var delta = PHR_FAST_DELTA

func BenchmarkScene(b *testing.B) {
	var bvh pt.BVH

	// Benchmark BVH construction speed
	b.Run("Build "+demoScene.Name, func(b *testing.B) {
		primitives := demoScene.Scene.UntransformedTracables()
		aux := pt.DefaultLBVH(primitives)
		builder := pt.NewPHRBuilder(primitives, alpha, delta, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})

	// Benchmark Trace speed for all view points
	buff := pt.NewBufferAspect(FRAME_SIZE, AR)
	for i, camera := range demoScene.Cameras {
		b.Run("Render view "+strconv.Itoa(i), func(b *testing.B) {
			renderer := pt.NewBenchmarkRenderer(bvh, camera)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderer.RenderToBuffer(buff)
			}
		})

		// Create image of render for validation
		img := buff.ToImage()
		f, err := os.Create(demoScene.Name + " " + strconv.Itoa(i) + ".png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, img)
	}
}

/*
func TestGridSearch(t *testing.T) {
	optimizer := pt.NewDefaultGridOptimizer(FRAME_SIZE, FRAME_SIZE)
	primitives := demoScene.Scene.Tracables()
	aux := pt.DefaultLBVH(primitives)
	a, d := optimizer.OptimizedPHRparams(aux, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
	fmt.Printf("Alpha: %v Delta: %v\n", a, d)
}


func TestBayOp(t *testing.T) {
	optimizer := pt.NewDefaultBayesianOptimizer(FRAME_SIZE, FRAME_SIZE)
	primitives := demoScene.Scene.Tracables()
	aux := pt.DefaultLBVH(primitives)
	a, d := optimizer.OptimizedPHRparams(aux, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
	fmt.Printf("Alpha: %v Delta: %v\n", a, d)
}
*/
