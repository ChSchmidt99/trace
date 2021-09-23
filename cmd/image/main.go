package main

import (
	"fmt"
	demo "github/chschmidt99/pt/pkg/demoscenes"
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
)

const (
	ASPECT_RATIO = 1
	FOV          = 60.0
	RESOLUTION   = 600
)

// Comment out which scenes should be rendered
var worlds = []demo.DemoScene{
	demo.CornellBox(),
	//demo.Bunny(),
	//demo.Dragon(),
	//demo.SponzaSun(),
	//demo.Buddha(),
	//demo.Hairball(),
	//demo.SibenikSun(),
	//demo.BreakfastSun(),
	//demo.FireplaceSun(),
	//demo.SanMiguelSun(),
}

func main() {
	for _, world := range worlds {
		camera := NewDefaultCamera(ASPECT_RATIO, FOV)
		bvh := world.Scene.CompilePHR(0.55, 9, 4)
		renderer := NewDefaultRenderer(bvh, camera)

		// Sky miss shader adds ambient light to all misses and adds a "sky" color interpolation
		renderer.Miss = SkyMissShader

		// Specify samples-per-pixel and max depth
		renderer.Spp = 100
		renderer.MaxDepth = 5
		renderer.Verbose = true

		// Render all view points
		for i, view := range world.ViewPoints {
			camera.SetTransformation(view)

			buff := NewPxlBufferAR(RESOLUTION, ASPECT_RATIO)
			renderer.RenderToBuffer(buff)

			// Write image to specified path
			path := fmt.Sprintf("%v_%v.png", world.Name, i)
			f, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			img := buff.ToImage()
			png.Encode(f, img)
			fmt.Printf("Written image to " + path + "\n")
		}
	}
}
