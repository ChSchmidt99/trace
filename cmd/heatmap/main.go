package main

import (
	"fmt"
	demo "github/chschmidt99/pt/pkg/demoscenes"
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
)

const (
	ASPECT_RATIO = 4.0 / 3
	FOV          = 70.0
	RESOLUTION   = 1200
)

func main() {
	// load any scene
	world := demo.Fireplace()

	// compile the scene using any provided BVH builder
	bvh := world.Scene.CompilePHR(0.55, 9, 2)
	//bvh := world.Scene.CompileLBVH()

	// specify at what traversal depth pixels start to be displayed in red
	threshold := 100

	// create a camera and set a view point
	camera := NewDefaultCamera(ASPECT_RATIO, FOV)
	camera.SetTransformation(world.ViewPoints[0])

	// create heatmap renderer, buffer and render to given buffer
	renderer := NewHeatMapRenderer(bvh, camera, threshold)
	buff := NewPxlBufferAR(RESOLUTION, ASPECT_RATIO)
	renderer.RenderToBuffer(buff)

	// write buffer to image in given directory
	path := world.Name + ".png"
	img := buff.ToImage()
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
	fmt.Printf("Written image to " + path + "\n")

}
