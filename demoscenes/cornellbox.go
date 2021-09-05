package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func CornellBox() *Scene {
	cube := ParseFromPath("../../assets/cube.obj")

	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	greenMat := Diffuse{Albedo: NewColor(.12, .45, .15)}
	redMat := Diffuse{Albedo: NewColor(.65, .05, .05)}
	lightMat := Light{Color: NewColor(5, 5, 5)}
	glass := Refractive{Albedo: NewColor(1, 1, 1), Ratio: 1.5}

	floor := NewSceneNode(NewMesh(cube, whiteMat))
	floor.Scale(10, 1, 10)

	back := NewSceneNode(NewMesh(cube, whiteMat))
	back.Scale(10, 10, 1)
	back.Translate(0, .5, 5)

	ceiling := NewSceneNode(NewMesh(cube, whiteMat))
	ceiling.Scale(11, 1, 11)
	ceiling.Translate(0, 10, 0)

	left := NewSceneNode(NewMesh(cube, greenMat))
	left.Scale(1, 10, 10)
	left.Translate(-5, .5, 0)

	right := NewSceneNode(NewMesh(cube, redMat))
	right.Scale(1, 10, 10)
	right.Translate(5, .5, 0)

	light := NewSceneNode(NewMesh(cube, lightMat))
	light.Scale(3, 1, 3)
	light.Translate(0, 9.99, 0)

	leftCube := NewSceneNode(NewMesh(cube, whiteMat))
	leftCube.Scale(2.5, 4.5, 2.5)
	leftCube.Translate(.8, .5, .8)
	leftCube.Rotate(NewVector3(0, 1, 0), 5)

	rightCube := NewSceneNode(NewMesh(cube, whiteMat))
	rightCube.Scale(3, 3, 3)
	rightCube.Translate(-.6, .6, 0)
	rightCube.Rotate(NewVector3(0, 1, 0), -5)

	sphere1 := Geometry{NewSphere(NewVector3(2, 2, 2), 1.5)}
	sphere2 := Geometry{NewSphere(NewVector3(-2, 2, -1), 1.5)}

	leftSphere := NewSceneNode(NewMesh(sphere1, glass))
	rightSphere := NewSceneNode(NewMesh(sphere2, glass))

	scene := NewScene()

	scene.Add(floor)
	scene.Add(back)
	scene.Add(ceiling)
	scene.Add(left)
	scene.Add(right)
	scene.Add(light)
	scene.Add(leftSphere)
	scene.Add(rightSphere)
	//scene.Add(leftCube)
	//scene.Add(rightCube)
	return scene
}
