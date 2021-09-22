package main

import (
	demo "github/chschmidt99/pt/pkg/demoscenes"
	. "github/chschmidt99/pt/pkg/pt"
	"strconv"
)

const (
	ASPECT_RATIO = 4.0 / 3
	FOV          = 55.0
	RESOLUTION   = 800
)

func main() {
	//world := demo.CornellBox()
	//world := demo.Bunny()
	//world := demo.Dragon()
	//world := demo.SanMiguel()
	//world := demo.Sponza()
	//world := demo.Buddha()
	//world := demo.Hairball()
	//world := demo.Sibenik()
	//world := demo.Breakfast()
	world := demo.Fireplace()

	camera := NewDefaultCamera(ASPECT_RATIO, FOV)
	//bvh := world.Scene.CompilePHR(0.65, 10, 4)
	bvh := world.Scene.CompileLBVH()
	renderer := NewDefaultRenderer(bvh, camera)
	//renderer.Closest = UnlitClosestHitShader
	//renderer.Miss = SkyMissShader
	//renderer.Miss = DawnMissShader
	renderer.Miss = SunMissShader
	renderer.Spp = 3000
	renderer.Verbose = true

	for i, view := range world.ViewPoints {
		camera.SetTransformation(view)
		renderer.RenderImageIncremental(world.Name+"_"+strconv.Itoa(i)+".png", RESOLUTION, ASPECT_RATIO, 10)
	}
}
