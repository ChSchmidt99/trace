package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Deer(ar, fov float64) DemoScene {
	geometry := ParseFromPath("../../assets/deer.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	bunny := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(bunny)

	camera := NewCamera(ar, fov, CameraTransformation{
		LookFrom: NewVector3(500, 250, 500),
		LookAt:   NewVector3(0, 110, 0),
		Up:       NewVector3(0, 1, 0),
	})

	return DemoScene{
		Name:    "Deer",
		Scene:   scene,
		Cameras: []*Camera{camera},
	}
}
