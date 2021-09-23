package main

import (
	demo "github/chschmidt99/pt/pkg/demoscenes"
	app "github/chschmidt99/pt/pkg/interactive"
	"github/chschmidt99/pt/pkg/pt"
	//. "github/chschmidt99/pt/pkg/pt"
)

const (
	ASPECT_RATIO = 1
	FOV          = 60.0
	RESOLUTION   = 300
)

/*
	Use w,a,s,d to move
	Press q to quit
*/
func main() {
	// select Scene to be rendered
	world := demo.Bunny()

	bvh := world.Scene.CompilePHR(0.55, 9, 4)
	camera := pt.NewCamera(ASPECT_RATIO, FOV, world.ViewPoints[0])
	renderer := pt.NewRealtimeRenderer(bvh, camera)
	cameraVelocity := 1.5
	runtime := app.NewInteractiveRuntime(renderer, ASPECT_RATIO, FOV, cameraVelocity, RESOLUTION)
	runtime.Run(nil)
}
