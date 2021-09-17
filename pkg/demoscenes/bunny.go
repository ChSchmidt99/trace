package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func Bunny(ar, fov float64) DemoScene {
	geometry := ParseFromPath("../../assets/local/bunny/bunny.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	//whiteMat := Refractive{Albedo: NewColor(.73, .73, .73), Ratio: 1.5}
	bunny := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(bunny)

	camera := NewCamera(ar, fov, CameraTransformation{
		LookFrom: NewVector3(0, 1, 2),
		LookAt:   NewVector3(-.25, .65, 0),
		Up:       NewVector3(0, 1, 0),
	})

	return DemoScene{
		Name:    "Bunny",
		Scene:   scene,
		Cameras: []*Camera{camera},
	}
}
