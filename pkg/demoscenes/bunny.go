package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Bunny() DemoScene {
	geometry := ParseFromPath("../../assets/local/bunny/bunny.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	bunny := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(bunny)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(0, 1, 2),
			LookAt:   NewVector3(-.25, .6, 0),
			Up:       NewVector3(0, 1, 0),
		},
		{
			LookFrom: NewVector3(-2, 2, 0),
			LookAt:   NewVector3(0, .5, 0),
			Up:       NewVector3(0, 1, 0),
		},
		{
			LookFrom: NewVector3(1.7, 1, 0),
			LookAt:   NewVector3(0, .7, 0),
			Up:       NewVector3(1, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Bunny",
		Scene:      scene,
		ViewPoints: views,
	}
}
