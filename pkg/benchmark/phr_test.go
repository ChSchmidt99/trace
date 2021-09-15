package benchmark

import (
	"fmt"
	"github/chschmidt99/pt/demoscenes"
	"github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
	"runtime"
	"testing"
)

const (
	AR               = 4.0 / 3.0
	FOV              = 50.0
	RESOLUTION       = 264
	PHR_FAST_ALPHA   = 0.5
	PHR_FAST_DELTA   = 6
	PHR_HQ_ALPHA     = 0.55
	PHR_HQ_DELTA     = 9
	BRANCHING_FACTOR = 2
)

var DEMO_SCENE *demoscenes.DemoScene

func loadDemoScene() demoscenes.DemoScene {
	if DEMO_SCENE == nil {
		demo := demoscenes.Bunny(AR, FOV)
		DEMO_SCENE = &demo
	}
	return *DEMO_SCENE
}

func BenchmarkPHR_Fast(b *testing.B) {
	scene := loadDemoScene().Scene
	camera := loadDemoScene().Cameras[0]
	name := loadDemoScene().Name

	var bvh pt.BVH
	b.Run("Build "+name, func(b *testing.B) {
		primitives := scene.Tracables()
		aux := pt.DefaultLBVH(primitives)
		builder := pt.NewPHRBuilder(primitives, PHR_FAST_ALPHA, PHR_FAST_DELTA, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})
	fmt.Printf("%v\n", bvh.Cost())
	buff := pt.NewBufferAspect(RESOLUTION, AR)
	b.Run("Render "+name, func(b *testing.B) {
		renderer := pt.NewBenchmarkRenderer(bvh, camera)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	img := buff.ToImage()
	f, err := os.Create("PHRFast_" + name + ".png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}

func BenchmarkPHR_HQ(b *testing.B) {
	scene := loadDemoScene().Scene
	camera := loadDemoScene().Cameras[0]
	name := loadDemoScene().Name

	var bvh pt.BVH
	b.Run("Build "+name, func(b *testing.B) {
		primitives := scene.Tracables()
		aux := pt.DefaultLBVH(primitives)
		builder := pt.NewPHRBuilder(primitives, PHR_HQ_ALPHA, PHR_HQ_DELTA, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})
	fmt.Printf("%v\n", bvh.Cost())
	buff := pt.NewBufferAspect(RESOLUTION, AR)
	b.Run("Render "+name, func(b *testing.B) {
		renderer := pt.NewBenchmarkRenderer(bvh, camera)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	img := buff.ToImage()
	f, err := os.Create("PHRHQ_" + name + ".png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}

func BenchmarkPHR_Grid(b *testing.B) {
	optimizer := pt.NewGridOptimizer([]float64{0.4, 0.45, 0.5, 0.55}, []int{5, 6, 7, 8})
	scene := loadDemoScene().Scene
	camera := loadDemoScene().Cameras[0]
	name := loadDemoScene().Name

	var bvh pt.BVH
	b.Run("Build "+name, func(b *testing.B) {
		primitives := scene.Tracables()
		aux := pt.DefaultLBVH(primitives)
		a, d := optimizer.OptimizedPHRparams(aux, camera, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		fmt.Printf("A: %v D:%v\n", a, d)
		builder := pt.NewPHRBuilder(primitives, a, d, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = builder.BuildFromAuxilary(aux)
		}
	})
	b.Run("Render "+name, func(b *testing.B) {
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
}

func BenchmarkRenderTime(b *testing.B) {
	scene := loadDemoScene().Scene
	camera := loadDemoScene().Cameras[0]
	primitives := scene.Tracables()
	aux := pt.DefaultLBVH(primitives)
	alpha := 0.55
	b.Run("5", func(b *testing.B) {
		builder := pt.NewPHRBuilder(primitives, alpha, 5, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		bvh := builder.BuildFromAuxilary(aux)
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	b.Run("6", func(b *testing.B) {
		builder := pt.NewPHRBuilder(primitives, alpha, 6, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		bvh := builder.BuildFromAuxilary(aux)
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	b.Run("7", func(b *testing.B) {
		builder := pt.NewPHRBuilder(primitives, alpha, 7, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		bvh := builder.BuildFromAuxilary(aux)
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	b.Run("8", func(b *testing.B) {
		builder := pt.NewPHRBuilder(primitives, alpha, 8, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		bvh := builder.BuildFromAuxilary(aux)
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	b.Run("9", func(b *testing.B) {
		builder := pt.NewPHRBuilder(primitives, alpha, 9, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		bvh := builder.BuildFromAuxilary(aux)
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	b.Run("10", func(b *testing.B) {
		builder := pt.NewPHRBuilder(primitives, alpha, 10, BRANCHING_FACTOR, runtime.GOMAXPROCS(0))
		bvh := builder.BuildFromAuxilary(aux)
		renderer := pt.NewDefaultRenderer(bvh, camera)
		renderer.Spp = 1
		buff := pt.NewBufferAspect(RESOLUTION, AR)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
}

func BenchmarkLBVH(b *testing.B) {
	scene := loadDemoScene().Scene
	camera := loadDemoScene().Cameras[0]
	name := loadDemoScene().Name

	var bvh pt.BVH
	b.Run("Build "+name, func(b *testing.B) {
		primitives := scene.Tracables()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bvh = pt.DefaultLBVH(primitives)
		}
	})
	fmt.Printf("%v\n", bvh.Cost())
	buff := pt.NewBufferAspect(RESOLUTION, AR)
	b.Run("Render "+name, func(b *testing.B) {
		renderer := pt.NewBenchmarkRenderer(bvh, camera)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			renderer.RenderToBuffer(buff)
		}
	})
	img := buff.ToImage()
	f, err := os.Create("LBVH_" + name + ".png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
