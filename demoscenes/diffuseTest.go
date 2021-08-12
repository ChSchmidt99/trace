package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func DiffuseTest() *Scene {

	whiteMat := Diffuse{Albedo: NewColor(.8, .8, 0)}
	greenMat := Diffuse{Albedo: NewColor(.7, .3, .3)}
	//whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	//greenMat := Diffuse{Albedo: NewColor(.12, .45, .15)}
	redMat := Diffuse{Albedo: NewColor(.65, .05, .05)}
	//lightMat := Light{Color: NewColor(1, 1, 1)}

	sphere2 := Geometry{NewSphere(NewVector3(-2, 2, -1), 2)}
	sphere1 := Geometry{NewSphere(NewVector3(2, 2, 2), 2)}
	floor := Geometry{NewSphere(NewVector3(0, -100, 0), 100)}

	leftSphere := NewSceneNode(NewMesh(sphere1, redMat))
	rightSphere := NewSceneNode(NewMesh(sphere2, greenMat))
	floorNode := NewSceneNode(NewMesh(floor, whiteMat))

	scene := NewScene()
	scene.Add(leftSphere)
	scene.Add(rightSphere)
	scene.Add(floorNode)
	return scene
}
