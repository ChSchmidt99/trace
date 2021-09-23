package demoscenes

import . "github/chschmidt99/pt/pkg/pt"

func SanMiguelSun() DemoScene {
	geometry := ParseFromPath("../../assets/local/san_miguel/san-miguel-low-poly.obj")

	//geometry := ParseFromPath("../../assets/local/san_miguel/san-miguel.obj")
	whiteMat := Diffuse{Albedo: NewColor(.83, .83, .83)}

	root := NewSceneNode(NewMesh(geometry, whiteMat))
	scene := NewScene()
	scene.Add(root)

	sphere := NewSphere(NewVector3(15, 25, 5), 6)
	light := Light{Color: NewColor(10, 6.5, 3)}
	sun := NewSceneNode(NewMesh(Geometry{sphere}, light))
	scene.Add(sun)

	views := []CameraTransformation{
		/*
			{
				LookFrom: NewVector3(22, 2, 7),
				LookAt:   NewVector3(15, 3, 1.5),
				Up:       NewVector3(0, 1, 0),
			},
		*/
		{
			LookFrom: NewVector3(15, 16, 10),
			LookAt:   NewVector3(15, 3, 1.5),
			Up:       NewVector3(0, 1, 0),
		},
		{
			LookFrom: NewVector3(14, 2, 9),
			LookAt:   NewVector3(15, 2, 7),
			Up:       NewVector3(0, 1, 0),
		},
		/*
			{
				LookFrom: NewVector3(22, 4, 5),
				LookAt:   NewVector3(5, 4, 5),
				Up:       NewVector3(0, 1, 0),
			},
		*/
		{
			LookFrom: NewVector3(26, 7, -2),
			LookAt:   NewVector3(5, 7, -2),
			Up:       NewVector3(0, 1, 0),
		},
		/*
			{
				LookFrom: NewVector3(20, 2, 8),
				LookAt:   NewVector3(0, 2, 0),
				Up:       NewVector3(0, 1, 0),
			},
		*/
	}
	return DemoScene{
		Name:       "San Miguel",
		Scene:      scene,
		ViewPoints: views,
	}
}
