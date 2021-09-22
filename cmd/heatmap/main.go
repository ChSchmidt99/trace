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
	ASPECT_RATIO = 4.0 / 3
	FOV          = 50.0
	RESOLUTION   = 400
)

func main() {
	//world := demo.CornellBox()
	//world := demo.Bunny()
	world := demo.SanMiguel()
	//world := demo.Hairball()

	camera := NewDefaultCamera(ASPECT_RATIO, FOV)
	bvh := world.Scene.CompilePHR(0.5, 6, 2)
	renderer := NewHeatMapRenderer(bvh, camera, 300)

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
