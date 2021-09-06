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
	RESOLUTION   = 1080
)

func main() {
	camera := NewCamera(ASPECT_RATIO, FOV, CameraTransformation{
		LookFrom: NewVector3(0, 5, -12),
		LookAt:   NewVector3(0, 5, 5),
		Up:       NewVector3(0, 1, 0),
	})

	scene := demo.CornellBox()
	bvh := scene.Compile()
	renderer := NewDefaultRenderer(bvh, camera)
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
