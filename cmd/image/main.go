package main

import (
	. "github/chschmidt99/pt/pkg/pt"
	"image/png"
	"os"
)

const (
	ASPECT_RATIO = 16.0 / 9
	FOV          = 70.0
	RESOLUTION   = 400
)

func main() {

	camera := NewCamera(ASPECT_RATIO, FOV, CameraTransformation{
		//LookFrom: NewVector3(7, 6, 10),
		LookFrom: NewVector3(10, 5, 10),
		LookAt:   NewVector3(0, 4, 0),
		Up:       NewVector3(0, 1, 0),
	})
	scene := NewScene()
	//cube := ParseFromPath("../../assets/cube.obj")
	//sphere := Geometry{NewSphere(NewVector3(15, 25, -40), 25)}

	triangle := Geometry{NewTriangleWithoutNormals(NewVector3(-25, 0, 40), NewVector3(-25, 0, -40), NewVector3(-25, 40, 40)),
		NewTriangleWithoutNormals(NewVector3(-25, 0, -40), NewVector3(-25, 40, -40), NewVector3(-25, 40, 40))}

	bunny := ParseFromPath("../../assets/local/bunny/bunny.obj")
	//tree := ParseFromPath("../../assets/tree.obj")

	mirrorMat := Reflective{Albedo: NewColor(.2, .2, .2), Diffusion: 0}
	greenMat := Diffuse{Albedo: NewColor255(47, 243, 84)}
	//glass := Refractive{Albedo: NewColor(.3, .3, .3), Ratio: 1.2}

	//diffuseTree := NewSceneNode(NewMesh(tree, &glass))
	//mirrorSphere := NewSceneNode(NewMesh(sphere, &mirrorMat))
	mirrorTriangle := NewSceneNode(NewMesh(triangle, &mirrorMat))
	bunnyNode := NewSceneNode(NewMesh(bunny, &greenMat))
	bunnyNode.ScaleUniform(50)

	//mirrorCube.ScaleUniform(50)
	//mirrorCube.Translate(0, 0, -2)

	scene.Add(bunnyNode)
	//scene.Add(diffuseTree)
	//scene.Add(mirrorSphere)
	scene.Add(mirrorTriangle)

	bvh := scene.Compile()
	renderer := NewDefaultRenderer(bvh, camera)
	buff := NewBufferAspect(RESOLUTION, ASPECT_RATIO)
	//renderer.RenderToBuffer(buff)
	renderer.IntersectionHeatMapToBuffer(buff)
	img := buff.ToImage()
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}
