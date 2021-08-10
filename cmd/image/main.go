package main

import (
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
)

func main() {
	ar := 16.0 / 9
	fov := 70.0
	camera := NewCamera(ar, fov, CameraTransformation{
		LookFrom: NewVector3(6, 6, 10),
		LookAt:   NewVector3(0, 6, 0),
		Up:       NewVector3(0, 1, 0),
	})
	scene := NewScene()

	/*
		radius := .5
		scene.Add(NewSceneNode(NewSphereMesh(NewVector3(2, 2, 1), radius, &Diffuse{
			Albedo: NewColor(1, 0, 0),
		})))
		scene.Add(NewSceneNode(NewTriangleMesh(NewVector3(0, 0, 2), NewVector3(2, 0, 0), NewVector3(0, 2, 0), &Reflective{
			Albedo:    NewColor(0, 0, 1),
			Diffusion: 0,
		})))
	*/

	cube := NewSceneNode(ParseFromPath("../../assets/deer.obj", &Diffuse{
		Albedo: NewColor(.5, .5, .5),
	}))
	cube.ScaleUniform(0.05)
	scene.Add(cube)

	bvh := scene.Compile()
	renderer := NewDefaultRenderer(bvh, camera)
	buff := NewBufferAspect(200, ar)
	renderer.RenderToBuffer(buff)
	img := buff.ToImage()
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
