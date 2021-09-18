package benchmark

import (
	"fmt"
	"github/chschmidt99/pt/pkg/demoscenes"
	"github/chschmidt99/pt/pkg/pt"
	"runtime"
	"strconv"
	"testing"
)

const (
	AR  = 1.0
	FOV = 60.0
	//FRAME_SIZE = 256
	FRAME_SIZE = 512
)

func BenchmarkBranchingFactor(b *testing.B) {
	world := demoscenes.Bunny()
	for i := 2; i <= 16; i *= 2 {
		bvh := BenchBuildPRH(b, "build_branch_"+strconv.Itoa(i), i, 0.4, 6, world)
		for j, view := range world.ViewPoints {
			BenchRender(b, "render_branch_"+strconv.Itoa(i)+"_view_"+strconv.Itoa(j), bvh, view)
		}
	}
}

func BenchmarkLBVH(b *testing.B) {
	world := demoscenes.Bunny()
	bvh := BenchBuildLBVH(b, "build_"+world.Name+"_LBVH", world)
	for i, view := range world.ViewPoints {
		BenchRender(b, "render_view_"+strconv.Itoa(i), bvh, view)
	}
}

func BenchmarkGridSearch(b *testing.B) {
	world := demoscenes.Bunny()
	optimizer := pt.NewDefaultGridOptimizer(FRAME_SIZE, FRAME_SIZE)
	primitives := world.Scene.UntransformedTracables()
	aux := pt.DefaultLBVH(primitives)
	var a float64
	var d float64
	b.Run("optimization", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a, d = optimizer.OptimizedPHRparams(aux, 2, runtime.GOMAXPROCS(0))
		}
	})
	bvh := BenchBuildPRH(b, "build", 2, a, d, world)
	for i, view := range world.ViewPoints {
		BenchRender(b, "render_view_"+strconv.Itoa(i), bvh, view)
	}
}

func BenchmarkBayOp(b *testing.B) {
	world := demoscenes.Bunny()
	optimizer := pt.NewDefaultBayesianOptimizer(FRAME_SIZE, FRAME_SIZE)
	primitives := world.Scene.UntransformedTracables()
	aux := pt.DefaultLBVH(primitives)
	var a float64
	var d float64
	b.Run("optimization", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a, d = optimizer.OptimizedPHRparams(aux, 2, runtime.GOMAXPROCS(0))
		}
	})

	var bvh pt.BVH
	b.Run("build", func(b *testing.B) {
		a, d = optimizer.OptimizedPHRparams(aux, 2, runtime.GOMAXPROCS(0))
		builder := pt.NewPHRBuilder(primitives, a, d, 2, runtime.GOMAXPROCS(0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})
	for i, view := range world.ViewPoints {
		buff := pt.NewFrameBufferAR(FRAME_SIZE, AR)
		camera := pt.NewCamera(AR, FOV, view)
		renderer := pt.NewBenchmarkRenderer(bvh, camera)

		b.Run("render_"+strconv.Itoa(i), func(b *testing.B) {
			a, d = optimizer.OptimizedPHRparams(aux, 2, runtime.GOMAXPROCS(0))
			builder := pt.NewPHRBuilder(primitives, a, d, 2, runtime.GOMAXPROCS(0))
			bvh = builder.BuildFromAuxilary(aux)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderer.RenderToBuffer(buff)
			}
		})
	}
}

func BenchmarkScene(b *testing.B) {

	/*
		primitives := world.Scene.UntransformedTracables()
		aux := pt.DefaultLBVH(primitives)
		builder := pt.NewPHRBuilder(primitives, 0, 0, 2, runtime.GOMAXPROCS(0))
		for alpha := 0.4; alpha < 0.7; alpha += 0.05 {
			for delta := 5; delta <= 10; delta += 1 {
					//buildName := fmt.Sprintf("build_%v_a_%v_d_%v", world.Name, 0.4, delta)
					//BenchBuildPRH(b, buildName, 2, 0.4, float64(delta), world)
					view := 2
					renderName := fmt.Sprintf("render_%v_%v_a_%v_d_%v", world.Name, view, alpha, delta)
					builder.Alpha = alpha
					builder.Delta = float64(delta)
					bvh := builder.BuildFromAuxilary(aux)
					BenchRender(b, renderName, bvh, world.ViewPoints[view])
			}
		}
	*/
}

func BenchPHR(b *testing.B, branching int, alpha, delta float64, world demoscenes.DemoScene) {
	buildName := fmt.Sprintf("build_%v_a_%v_d_%v", world.Name, alpha, delta)
	bvh := BenchBuildPRH(b, buildName, branching, alpha, delta, world)
	for i, view := range world.ViewPoints {
		renderName := fmt.Sprintf("render_view_%v_a_%v_d_%v", i, alpha, delta)
		BenchRender(b, renderName, bvh, view)
	}
}

func BenchBuildPRH(b *testing.B, name string, branching int, alpha, delta float64, world demoscenes.DemoScene) pt.BVH {
	var bvh pt.BVH
	primitives := world.Scene.UntransformedTracables()
	aux := pt.DefaultLBVH(primitives)
	builder := pt.NewPHRBuilder(primitives, alpha, delta, branching, runtime.GOMAXPROCS(0))
	b.Run(name, func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})
	return bvh
}

func BenchBuildLBVH(b *testing.B, name string, world demoscenes.DemoScene) pt.BVH {
	var bvh pt.BVH
	primitives := world.Scene.UntransformedTracables()
	b.Run(name, func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = pt.DefaultLBVH(primitives)
		}
	})
	return bvh
}

func BenchRender(b *testing.B, name string, bvh pt.BVH, view pt.CameraTransformation) {
	// Benchmark Trace speed for all view points
	buff := pt.NewFrameBufferAR(FRAME_SIZE, AR)
	camera := pt.NewCamera(AR, FOV, view)
	renderer := pt.NewBenchmarkRenderer(bvh, camera)
	b.Run(name, func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
}

func TestPHR(t *testing.T) {
	world := demoscenes.Bunny()
	primitives := world.Scene.UntransformedTracables()
	aux := pt.DefaultLBVH(primitives)
	fmt.Printf("%v\n", aux.Cost())
	/*
		builder := pt.NewPHRBuilder(primitives, 0, 0, 2, runtime.GOMAXPROCS(0))
		for alpha := 0.4; alpha < 0.7; alpha += 0.05 {
			for delta := 5; delta <= 10; delta += 1 {
				builder.Alpha = alpha
				builder.Delta = float64(delta)
				_, cost := builder.BuildWithCost(aux)
				fmt.Printf("%v\n", cost)
			}
		}
	*/
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
