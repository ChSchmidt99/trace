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
	RESOLUTION   = 800
)

func main() {
	//demoScene := demo.CornellBox(ASPECT_RATIO, FOV)
	demoScene := demo.Bunny(ASPECT_RATIO, FOV)
	//demoScene := demo.SanMiguel(ASPECT_RATIO, FOV)

	bvh := demoScene.Scene.CompilePHR(0.5, 6, 2)
	//bvh := demoScene.Scene.CompileLBVH()

	renderer := NewHeatMapRenderer(bvh, demoScene.Cameras[0])
	buff := NewBufferAspect(RESOLUTION, ASPECT_RATIO)
	renderer.RenderToBuffer(buff)
	img := buff.ToImage()
	f, err := os.Create(demoScene.Name + ".png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
