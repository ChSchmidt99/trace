package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func FireplaceSun() DemoScene {
	geometry := ParseFromPath("../../assets/local/fireplace/fireplace_room_window.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)
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
