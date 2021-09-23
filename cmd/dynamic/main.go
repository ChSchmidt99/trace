package main

import (
	app "github/chschmidt99/pt/pkg/interactive"
	"github/chschmidt99/pt/pkg/pt"
	"time"
)

const (
	ASPECT_RATIO = 1
	FOV          = 60.0
	RESOLUTION   = 300
)

// Very rudimentary example of a dynamic scene
var scene *pt.Scene
var renderer pt.Renderer

/*
	Use w,a,s,d to move
	Press q to quit
*/
func main() {
	// Create a demo scene
	scene = pt.NewScene()
	whiteMat := pt.Diffuse{Albedo: pt.NewColor(.73, .73, .73)}
	sphere := pt.NewMesh(pt.Geometry{pt.NewSphere(pt.NewVector3(5, 2, 0), 2)}, whiteMat)
	geometry := pt.ParseFromPath("../../assets/cube.obj")
	floor := pt.NewSceneNode(pt.NewMesh(geometry, whiteMat))
	floor.Scale(100, 1, 10)
	floor.Translate(0, -1, 0)
	scene.Add(floor)
	sphere1 := pt.NewSceneNode(sphere)
	sphere2 := pt.NewSceneNode(sphere)
	sphere3 := pt.NewSceneNode(sphere)
	sphere4 := pt.NewSceneNode(sphere)
	sphere5 := pt.NewSceneNode(sphere)
	scene.Add(sphere2)
	scene.Add(sphere3)
	scene.Add(sphere4)
	scene.Add(sphere5)
	scene.Add(sphere1)
	centerSphereRadius := 5.0
	blueMat := pt.Reflective{Albedo: pt.NewColor(0.2, 0.2, 0.8), Diffusion: 0}
	centerMesh := pt.NewMesh(pt.Geometry{pt.NewSphere(pt.NewVector3(0, centerSphereRadius, 0), centerSphereRadius)}, &blueMat)
	centerSphere := pt.NewSceneNode(centerMesh)
	scene.Add(centerSphere)

	// Create a few random animations
	wps := []pt.Vector3{
		pt.NewVector3(-20, 0, -20),
		pt.NewVector3(20, 0, -20),
		pt.NewVector3(20, 0, 20),
		pt.NewVector3(-20, 0, 20),
		pt.NewVector3(-20, 0, -20),
	}
	path := app.NewUniformSequence(wps, 0, time.Second*10)
	animation := app.NewAnimation(sphere1, path, time.Now(), true)

	wps2 := []pt.Vector3{
		pt.NewVector3(-20, -20, -20),
		pt.NewVector3(20, 20, 20),
		pt.NewVector3(-20, -20, -20),
	}
	path2 := app.NewUniformSequence(wps2, 0, time.Second*5)
	animation2 := app.NewAnimation(sphere2, path2, time.Now(), true)

	wps3 := []pt.Vector3{
		pt.NewVector3(0, -20, 0),
		pt.NewVector3(0, 20, 0),
		pt.NewVector3(0, -20, 0),
	}
	path3 := app.NewUniformSequence(wps3, 0, time.Second*15)
	animation3 := app.NewAnimation(sphere3, path3, time.Now(), true)

	wps4 := []pt.Vector3{
		pt.NewVector3(-20, -20, 20),
		pt.NewVector3(20, 20, -20),
		pt.NewVector3(-20, -20, 20),
	}
	path4 := app.NewUniformSequence(wps4, 0, time.Second*15)
	animation4 := app.NewAnimation(sphere4, path4, time.Now(), true)

	wps5 := []pt.Vector3{
		pt.NewVector3(-10, -10, -10),
		pt.NewVector3(10, 10, -10),
		pt.NewVector3(10, -10, 10),
		pt.NewVector3(-10, 10, 10),
		pt.NewVector3(-10, -10, -10),
	}
	path5 := app.NewUniformSequence(wps5, 0, time.Second*15)
	animation5 := app.NewAnimation(sphere5, path5, time.Now(), true)

	bvh := scene.CompileLBVH()
	camera := pt.NewCamera(ASPECT_RATIO, FOV, pt.CameraTransformation{
		LookFrom: pt.NewVector3(15, 15, 15),
		LookAt:   pt.NewVector3(0, 0, 0),
		Up:       pt.NewVector3(0, 1, 0),
	})

	// Execute the renderer
	renderer = pt.NewRealtimeRenderer(bvh, camera)
	cameraVelocity := 1.5
	runtime := app.NewInteractiveRuntime(renderer, ASPECT_RATIO, FOV, cameraVelocity, RESOLUTION)
	runtime.AddAnimation(animation)
	runtime.AddAnimation(animation2)
	runtime.AddAnimation(animation3)
	runtime.AddAnimation(animation4)
	runtime.AddAnimation(animation5)
	runtime.Run(Update)
}

func Update() {
	// Executet after every scene change
	bvh := scene.CompileLBVH()
	renderer.SetBvh(bvh)
}
