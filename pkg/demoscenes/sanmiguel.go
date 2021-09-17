package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func SanMiguel(ar, fov float64) DemoScene {
	geometry := ParseFromPath("../../assets/local/san_miguel/san-miguel-low-poly.obj")
	//geometry := ParseFromPath("../../assets/local/san_miguel/san-miguel.obj")
	whiteMat := Diffuse{Albedo: NewColor(.73, .73, .73)}
	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	camera := NewCamera(ar, fov, CameraTransformation{
		LookFrom: NewVector3(22, 2, 7),
		LookAt:   NewVector3(15, 3, 1.5),
		Up:       NewVector3(0, 1, 0),
	})

	return DemoScene{
		Name:    "San Miguel",
		Scene:   scene,
		Cameras: []*Camera{camera},
	}
}
