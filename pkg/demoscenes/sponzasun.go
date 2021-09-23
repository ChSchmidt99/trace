package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func SponzaSun() DemoScene {
	geometry := ParseFromPath("../../assets/local/sponza/sponza.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}

	sponza := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(sponza)

	sphere := NewSphere(NewVector3(-10, 15, 0), 4)
	light := Light{Color: NewColor(10, 6.5, 3)}
	sun := NewSceneNode(NewMesh(Geometry{sphere}, light))
	scene.Add(sun)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(-5, 3, 0),
			LookAt:   NewVector3(5, 3, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Sponza",
		Scene:      scene,
		ViewPoints: views,
	}
}
