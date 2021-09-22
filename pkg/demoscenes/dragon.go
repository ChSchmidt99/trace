package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Dragon() DemoScene {
	geometry := ParseFromPath("../../assets/local/dragon/dragon.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	//whiteMat := Refractive{Albedo: NewColor(.73, .73, .73), Ratio: 1.5}
	bunny := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(bunny)

	views := []CameraTransformation{
		{
			LookFrom: NewVector3(-.75, .75, -.3),
			LookAt:   NewVector3(.1, 0, 0),
			Up:       NewVector3(0, 1, 0),
		},
		{
			LookFrom: NewVector3(-1, .25, 0),
			LookAt:   NewVector3(0, .1, 0),
			Up:       NewVector3(0, 1, 0),
		},
		{
			LookFrom: NewVector3(0, .25, -1),
			LookAt:   NewVector3(0, 0, 0),
			Up:       NewVector3(0, 1, 0),
		},
	}

	return DemoScene{
		Name:       "Dragon",
		Scene:      scene,
		ViewPoints: views,
	}
}
