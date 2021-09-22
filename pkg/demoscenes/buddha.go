package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Buddha() DemoScene {
	geometry := ParseFromPath("../../assets/local/buddha/buddha.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(0, 0, -1),
			LookAt:   NewVector3(0, 0, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Buddha",
		Scene:      scene,
		ViewPoints: views,
	}
}
