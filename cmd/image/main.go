package main

import (
	"fmt"
	demo "github/chschmidt99/pt/pkg/demoscenes"
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
)

const (
	ASPECT_RATIO = 4.0 / 3.0
	FOV          = 70.0
	RESOLUTION   = 1200
)

var worlds = []demo.DemoScene{
	//demo.CornellBox(),
	//demo.Bunny(),
	//demo.Dragon(),
	//demo.Sponza(),
	//demo.Buddha(),
	//demo.Hairball(),
	//demo.Sibenik(),
	//demo.Breakfast(),
	demo.Fireplace(),
	//demo.SanMiguelSun(),
}

func main() {
	//world := demo.CornellBox()
	//world := demo.Bunny()
	//world := demo.Dragon()
	//world := demo.SanMiguel()
	//world := demo.Sponza()
	//world := demo.Buddha()
	//world := demo.Hairball()
	//world := demo.Sibenik()
	//world := demo.Breakfast()
	//world := demo.Fireplace()

	for _, world := range worlds {
		camera := NewDefaultCamera(ASPECT_RATIO, FOV)
		bvh := world.Scene.CompileLBVH()
		renderer := NewDefaultRenderer(bvh, camera)
		renderer.Miss = SkyMissShader
		renderer.Spp = 1000
		renderer.MaxDepth = 10
		renderer.Verbose = true
		for i, view := range world.ViewPoints {
			camera.SetTransformation(view)
			buff := NewPxlBufferAR(RESOLUTION, ASPECT_RATIO)
			renderer.RenderToBuffer(buff)
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
