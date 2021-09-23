package benchmark

import (
	"fmt"
	"github/chschmidt99/pt/pkg/demoscenes"
	"github/chschmidt99/pt/pkg/pt"
	"runtime"
	"testing"
)

const (
	AR             = 1.0
	FOV            = 60.0
	PHR_FAST_ALPHA = 0.5
	PHR_FAST_DELTA = 6
	PHR_HQ_ALPHA   = 0.55
	PHR_HQ_DELTA   = 9
)

func BenchmarkThreads(b *testing.B) {
	world := demoscenes.Bunny()
	builder := pt.NewPHRBuilder(world.Scene.UntransformedTracables(), 0.5, 6, 4, runtime.GOMAXPROCS(0))
	bvh := builder.Build()
	benchRender(b, "", bvh, world, 256)
}

func BenchmarkScene(b *testing.B) {
	// Select Scene to be benchmarked
	world := demoscenes.Fireplace()

	// Specify all tested resolutions
	resolutions := []int{256, 512}

	// Specify branching factor of BVH
	branching := 2
	benchLBVH(b, world, resolutions)
	benchFullPHR(b, "PHR_Fast", PHR_FAST_ALPHA, PHR_FAST_DELTA, branching, world, resolutions)
	benchFullPHR(b, "PHR_HQ", PHR_HQ_ALPHA, PHR_HQ_DELTA, branching, world, resolutions)
	benchGridSearch(b, branching, world, resolutions)
	benchBayOp(b, branching, world, resolutions)
}

func benchLBVH(b *testing.B, world demoscenes.DemoScene, resolutions []int) {
	bvh := benchBuildLBVH(b, world)
	for _, resolution := range resolutions {
		benchRender(b, "lbvh", bvh, world, resolution)
	}
}

func benchGridSearch(b *testing.B, branching int, world demoscenes.DemoScene, resolutions []int) {
	name := "Grid_Search"
	for _, res := range resolutions {
		optimizer := pt.NewDefaultGridOptimizer(res, res)
		primitives := world.Scene.UntransformedTracables()
		aux := pt.DefaultLBVH(primitives)
		var a float64
		var d float64
		b.Run(name+"_"+world.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a, d = optimizer.OptimizedPHRparams(aux, branching, runtime.GOMAXPROCS(0))
			}
		})
		bvh := benchBuildPHR(b, name, branching, a, d, world)
		benchRender(b, name, bvh, world, res)
	}
}

func benchBayOp(b *testing.B, branching int, world demoscenes.DemoScene, resolutions []int) {
	name := "Bayesian_Optimization"
	for _, res := range resolutions {
		optimizer := pt.NewDefaultBayesianOptimizer(res, res)
		primitives := world.Scene.UntransformedTracables()
		aux := pt.DefaultLBVH(primitives)
		var a float64
		var d float64
		b.Run(name+" "+world.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a, d = optimizer.OptimizedPHRparams(aux, branching, runtime.GOMAXPROCS(0))
			}
		})
		bvh := benchBuildPHR(b, name, branching, a, d, world)
		benchRender(b, name, bvh, world, res)
	}
}

func benchFullPHR(b *testing.B, name string, alpha, delta float64, branching int, world demoscenes.DemoScene, resolutions []int) {
	bvh := benchBuildPHR(b, name, branching, alpha, delta, world)
	for _, resolution := range resolutions {
		benchRender(b, name, bvh, world, resolution)
	}
}

func benchBuildPHR(b *testing.B, name string, branching int, alpha, delta float64, world demoscenes.DemoScene) pt.BVH {
	var bvh pt.BVH
	primitives := world.Scene.UntransformedTracables()
	aux := pt.DefaultLBVH(primitives)
	builder := pt.NewPHRBuilder(primitives, alpha, delta, branching, runtime.GOMAXPROCS(0))
	n := fmt.Sprintf("%v_build_phr_a_%.4f_d_%.4f_%v", name, alpha, delta, world.Name)
	b.Run(n, func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})
	return bvh
}

func benchBuildLBVH(b *testing.B, world demoscenes.DemoScene) pt.BVH {
	var bvh pt.BVH
	primitives := world.Scene.UntransformedTracables()
	b.Run("build_lbvh_"+world.Name, func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = pt.DefaultLBVH(primitives)
		}
	})
	return bvh
}

func benchRender(b *testing.B, name string, bvh pt.BVH, world demoscenes.DemoScene, frameSize int) {
	buff := pt.NewFrameBufferAR(frameSize, AR)
	camera := pt.NewDefaultCamera(AR, FOV)
	renderer := pt.NewBenchmarkRenderer(bvh, camera)
	//for viewIndex, view := range world.ViewPoints {
	view, viewIndex := world.ViewPoints[0], 0
	camera.SetTransformation(view)
	n := fmt.Sprintf(name+"_render_%v_view_%v_%vx%v", world.Name, viewIndex, frameSize, frameSize)
	b.Run(n, func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	//}

}
