package main

import (
	"fmt"
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
	//scene, camera := demo.CornellBox(ASPECT_RATIO, FOV)
	//scene, camera := demo.Bunny(ASPECT_RATIO, FOV)
	demoScene := demo.SanMiguel(ASPECT_RATIO, FOV)

	fmt.Printf("Prims: %v\n", len(demoScene.Scene.Tracables()))

	bvh := demoScene.Scene.Compile()

	renderer := NewHeatMapRenderer(bvh, demoScene.Cameras[0])
	buff := NewBufferAspect(RESOLUTION, ASPECT_RATIO)
	renderer.RenderToBuffer(buff)
	img := buff.ToImage()
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
