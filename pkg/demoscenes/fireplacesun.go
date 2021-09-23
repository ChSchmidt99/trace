package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Fireplace() DemoScene {
	geometry := ParseFromPath("../../assets/local/fireplace/fireplace_room_window.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	sphere := NewSphere(NewVector3(-30, 30, 30), 20)
	light := Light{Color: NewColor(8, 6, 5)}
	sun := NewSceneNode(NewMesh(Geometry{sphere}, light))
	scene.Add(sun)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(5, 2, -3),
			LookAt:   NewVector3(0, 0, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Fireplace",
		Scene:      scene,
		ViewPoints: views,
	}
}
