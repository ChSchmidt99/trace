package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Hairball(ar, fov float64) DemoScene {
	geometry := ParseFromPath("../../assets/local/hairball.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(7, 7, 7),
			LookAt:   NewVector3(0, 0, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Hairball",
		Scene:      scene,
		ViewPoints: views,
	}
}
