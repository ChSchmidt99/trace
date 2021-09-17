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
	RESOLUTION   = 800
)

func main() {
	//world := demo.CornellBox(ASPECT_RATIO, FOV)
	world := demo.Bunny()
	//world := demo.SanMiguel(ASPECT_RATIO, FOV)
	//world := demo.Hairball(ASPECT_RATIO, FOV)

	camera := NewDefaultCamera(ASPECT_RATIO, FOV)
	bvh := world.Scene.CompilePHR(0.5, 6, 2)
	renderer := NewHeatMapRenderer(bvh, camera)

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
		fmt.Printf("Written image to " + imageName + ".png\n")
	}
}
