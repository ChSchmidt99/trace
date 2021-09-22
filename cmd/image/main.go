package main

import (
	"fmt"
	demo "github/chschmidt99/pt/pkg/demoscenes"
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
	"strconv"
)

const (
	ASPECT_RATIO = 1
	FOV          = 60.0
	RESOLUTION   = 300
)

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
	world := demo.Fireplace()

	camera := NewDefaultCamera(ASPECT_RATIO, FOV)
	bvh := world.Scene.CompilePHR(0.55, 9, 4)
	renderer := NewDefaultRenderer(bvh, camera)
	//renderer.Closest = UnlitClosestHitShader
	renderer.Miss = SkyMissShader
	renderer.Spp = 300
	renderer.Verbose = true

	for i, view := range world.ViewPoints {
		camera.SetTransformation(view)
		buff := NewPxlBufferAR(RESOLUTION, ASPECT_RATIO)
		renderer.RenderToBuffer(buff)
		img := buff.ToImage()
		imageName := world.Name + " " + strconv.Itoa(i) + ".png"
		f, err := os.Create(imageName)
		if err != nil {
			panic(err)
		}
		png.Encode(f, img)
		fmt.Printf("Written image to " + imageName + "\n")
	}
}
