package main

import (
	. "github/chschmidt99/pt/pkg/pt"
	"image"
	"image/png"
	"os"
)

func main() {
	ar := 16.0 / 9
	fov := 60.0
	camera := NewCamera(ar, fov, CameraTransformation{
		LookFrom: NewVector3(0, 2, 2),
		LookAt:   NewVector3(0, 0, 0),
		Up:       NewVector3(0, 1, 0),
	})
	scene := NewScene()
	radius := .8
	scene.Add(NewSceneNode(NewSphereMesh(NewVector3(1, 0, 0), radius, &Diffuse{
		Albedo: NewColor(1, 0, 0),
	})))
	scene.Add(NewSceneNode(NewSphereMesh(NewVector3(-1, 0, 0), radius, &Reflective{
		Albedo:    NewColor(0, 0, 1),
		Diffusion: 0,
	})))
	bvh := scene.Compile()
	renderer := NewDefaultRenderer(bvh, camera)
	buff := NewBufferAspect(200, ar)
	renderer.RenderToBuffer(buff)
	img := toImage(buff)
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}

func toImage(buffer *PixelBuffer) image.Image {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{buffer.Width, buffer.Height}
	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})
	for i, color := range buffer.Buff {
		x := i % buffer.Width
		y := i / buffer.Width
		img.Set(x, y, color.GoColor())
	}
	return img
}
