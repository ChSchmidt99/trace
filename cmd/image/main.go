package main

import (
	demo "github/chschmidt99/pt/demoscenes"
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
)

const (
	ASPECT_RATIO = 4.0 / 3
	FOV          = 50.0
	RESOLUTION   = 200
)

func main() {
	//scene, camera := demo.CornellBox(ASPECT_RATIO, FOV)
	scene, camera := demo.Bunny(ASPECT_RATIO, FOV)

	bvh := scene.CompileLBVH()
	renderer := NewDefaultRenderer(bvh, camera)
	renderer.Closest = UnlitClosestHitShader
	renderer.Miss = SkyMissShader
	renderer.Spp = 200

	buff := NewBufferAspect(RESOLUTION, ASPECT_RATIO)
	renderer.RenderToBuffer(buff)
	//renderer.IntersectionHeatMapToBuffer(buff)
	img := buff.ToImage()
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
