package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func BreakfastSun() DemoScene {
	geometry := ParseFromPath("../../assets/local/breakfast/breakfast_room.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	sphere := NewSphere(NewVector3(5, 4, 5), 2)
	light := Light{Color: NewColor(8, 6, 5)}
	sun := NewSceneNode(NewMesh(Geometry{sphere}, light))
	scene.Add(sun)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(1, 3.5, 6),
			LookAt:   NewVector3(-.5, 1.5, -1),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Breakfast",
		Scene:      scene,
		ViewPoints: views,
	}
}
