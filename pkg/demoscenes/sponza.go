package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Sponza() DemoScene {
	geometry := ParseFromPath("../../assets/local/sponza/sponza.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}

	sponza := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(sponza)

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
