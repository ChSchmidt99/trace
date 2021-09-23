package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func SibenikSun() DemoScene {
	geometry := ParseFromPath("../../assets/local/sibenik/sibenik.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	sphere := NewSphere(NewVector3(0, 2, 0), 2)
	light := Light{Color: NewColor(10, 6.5, 3)}
	sun := NewSceneNode(NewMesh(Geometry{sphere}, light))
	scene.Add(sun)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(-16, -10, 0),
			LookAt:   NewVector3(1, -10, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Sibenik",
		Scene:      scene,
		ViewPoints: views,
	}
}
