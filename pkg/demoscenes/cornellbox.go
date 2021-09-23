package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func CornellBox() DemoScene {
	cube := ParseFromPath("../../assets/cube.obj")

	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	teal := Diffuse{Albedo: NewColor(0.07, .56, .77)}
	redMat := Diffuse{Albedo: NewColor(.77, .07, .21)}
	lightMat := Light{Color: NewColor(10, 10, 10)}
	mirror := Reflective{Albedo: NewColor(.25, .25, .25), Diffusion: 0.1}
	glass := Refractive{Albedo: NewColor(1, 1, 1), Ratio: 1.5}

	floor := NewSceneNode(NewMesh(cube, mirror))
	floor.Scale(10, 1, 20)

	back := NewSceneNode(NewMesh(cube, whiteMat))
	back.Scale(10, 10, 1)
	back.Translate(0, .5, 5)

	front := NewSceneNode(NewMesh(cube, whiteMat))
	front.Scale(10, 10, 1)
	front.Translate(0, .5, -10)

	ceiling := NewSceneNode(NewMesh(cube, whiteMat))
	ceiling.Scale(11, 1, 20)
	ceiling.Translate(0, 10, 0)

	left := NewSceneNode(NewMesh(cube, teal))
	left.Scale(1, 10, 20)
	left.Translate(-5, .5, 0)

	right := NewSceneNode(NewMesh(cube, redMat))
	right.Scale(1, 10, 20)
	right.Translate(5, .5, 0)

	light := NewSceneNode(NewMesh(cube, lightMat))
	light.Scale(3, 1, 3)
	light.Translate(0, 9.99, 0)

	leftCube := NewSceneNode(NewMesh(cube, whiteMat))
	leftCube.Scale(2.5, 4.5, 2.5)
	leftCube.Translate(.8, .5, .8)
	leftCube.Rotate(NewVector3(0, 1, 0), .45)

	rightCube := NewSceneNode(NewMesh(cube, whiteMat))
	rightCube.Scale(3, 3, 3)
	rightCube.Rotate(NewVector3(0, 1, 0), 45)
	rightCube.Translate(-.6, .6, 0)

	sphere2 := Geometry{NewSphere(NewVector3(-2, 2, -1), 1.5)}
	rightSphere := NewSceneNode(NewMesh(sphere2, glass))

	scene := NewScene()
	scene.Add(floor)
	scene.Add(back)
	scene.Add(front)
	scene.Add(ceiling)
	scene.Add(left)
	scene.Add(right)
	scene.Add(light)
	scene.Add(rightSphere)
	scene.Add(leftCube)
	views := []CameraTransformation{
		{
			LookFrom: NewVector3(0, 5, -9),
			LookAt:   NewVector3(0, 5, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Cornell Box",
		Scene:      scene,
		ViewPoints: views,
	}
}
