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
		//LookFrom: NewVector3(7, 6, 10),
		LookFrom: NewVector3(20, 20, 20),
		LookAt:   NewVector3(0, 10, 0),
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
	deerGeometry := ParseFromPath("../../assets/tree.obj")

	diffuseDeer := NewSceneNode(NewMesh(deerGeometry, &Diffuse{Albedo: NewColor(1, .5, .5)}))
	diffuseDeer.ScaleUniform(0.05)
	diffuseDeer.Translate(50, 0, 0)
	//scene.Add(diffuseDeer)

	reflectiveDeer := NewSceneNode(NewMesh(deerGeometry, &Reflective{Albedo: NewColor(.2, .2, .2), Diffusion: 0}))
	//reflectiveDeer := NewSceneNode(NewMesh(deerGeometry, &Diffuse{Albedo: NewColor(.2, .2, .2)}))
	//reflectiveDeer.ScaleUniform(0.045)
	//reflectiveDeer.Rotate(NewVector3(0, 1, 0), 45)
	//reflectiveDeer.Translate(-100, 0, -50)
	scene.Add(reflectiveDeer)

	bvh := scene.Compile()
	renderer := NewDefaultRenderer(bvh, camera)
	buff := NewBufferAspect(720, ar)
	renderer.RenderToBuffer(buff)
	img := buff.ToImage()
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
