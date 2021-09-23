package raytracer

import (
	"fmt"
	"github/chschmidt99/pt/pkg/pt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Disclaimer: interactive.go is experimental and not clean code!

type Runtime struct {
	renderer pt.Renderer

	frameCounter int
	pressedKeys  map[sdl.Keycode]time.Time // Holds the keycode as key and the timestamp of the press as value

	mouseEvents      []*MouseEvent
	pitch            float64
	yaw              float64
	cameraVelocity   float64
	mouseSensitivity float64
	lastFrame        time.Time
	window           *Window

	scheduler *animationScheduler
	onRender  func()
}

// scene: Scene to be rendered, camera has to be of type RealTimeCamera; BVH already calculated
// aspectRatio: window aspect ratio
// cameraVelocity: how fast the camera will move (translation, not pitch and yaw)
// resolution: height of source frameBuffer (720 for 720p, 1080 for 1080p etc.)
// scale: factor to scale up the window size for low resolutions
// maxThreads: number of threads for render parallelization
func NewInteractiveRuntime(renderer pt.Renderer, aspectRatio float64, fov float64, cameraVelocity float64, resolution int) *Runtime {
	return &Runtime{
		renderer:         renderer,
		pressedKeys:      make(map[sdl.Keycode]time.Time),
		mouseEvents:      make([]*MouseEvent, 0),
		pitch:            10,
		yaw:              45,
		mouseSensitivity: 0.05,
		lastFrame:        time.Now(),
		frameCounter:     0,
		scheduler:        newAnimationScheduler(),
		cameraVelocity:   cameraVelocity,
		window:           NewWindow(resolution, aspectRatio, "Press 'Q' to quit"),
	}
}

func (r *Runtime) Run(onRender func()) {
	r.onRender = onRender
	r.window.InitWindow()
	go r.runFPSprinter()
	r.window.Run(r.update, r.receiveEvent)
	r.window.Close()

}

func (r *Runtime) AddAnimation(a *Animation) {
	r.scheduler.queueAnimation(a)
}

func (runtime *Runtime) update(window *Window) {
	// translate camera
	deltaTime := time.Since(runtime.lastFrame)
	runtime.lastFrame = time.Now()
	velocity := runtime.cameraVelocity * float64(deltaTime.Seconds())

	cam := runtime.renderer.GetCamera()
	for keyCode := range runtime.pressedKeys {
		switch keyCode {
		case sdl.K_w:
			cam.Translate(cam.W().Mul(-velocity))
		case sdl.K_s:
			cam.Translate(cam.W().Mul(velocity))
		case sdl.K_a:

			dir := cam.W().Cross(cam.Up()).Unit()
			cam.Translate(dir.Mul(velocity))
		case sdl.K_d:
			dir := cam.W().Cross(cam.Up()).Unit()
			cam.Translate(dir.Mul(-velocity))
		}
	}

	// rotate camera
	xOffset := 0.0
	yOffset := 0.0
	for _, event := range runtime.mouseEvents {
		xOffset += float64(event.x)
		yOffset += float64(event.y)
	}
	runtime.mouseEvents = make([]*MouseEvent, 0)

	xOffset *= runtime.mouseSensitivity
	yOffset *= runtime.mouseSensitivity

	runtime.yaw += xOffset
	runtime.pitch += yOffset

	if runtime.pitch > 89 {
		runtime.pitch = 89
	} else if runtime.pitch < -89 {
		runtime.pitch = -89
	}

	x := math.Cos(degreesToRadians(runtime.yaw)) * math.Cos(degreesToRadians(runtime.pitch))
	y := math.Sin(degreesToRadians(runtime.pitch))
	z := math.Sin(degreesToRadians(runtime.yaw)) * math.Cos(degreesToRadians(runtime.pitch))
	cam.SetFront(pt.NewVector3(x, y, z))

	runtime.scheduler.animationTick(deltaTime)

	if runtime.onRender != nil {
		runtime.onRender()
	}
	runtime.renderer.RenderToBuffer(window.buffer)
	window.Refresh()
	runtime.frameCounter++
}

func (runtime *Runtime) receiveEvent(event interface{}) {
	switch typedEvent := event.(type) {
	case *KeyEvent:
		if typedEvent.keyCode == sdl.K_q {
			runtime.quit()
		}
		if typedEvent.eventType == DOWN {
			runtime.pressedKeys[typedEvent.keyCode] = typedEvent.timestamp
		} else {
			delete(runtime.pressedKeys, typedEvent.keyCode)
		}
	case *MouseEvent:
		runtime.mouseEvents = append(runtime.mouseEvents, typedEvent)
	}
}

// Simplyfied FPS printer to get a rough estimate
func (runtime *Runtime) runFPSprinter() {
	for {
		time.Sleep(time.Second)
		fmt.Println("FPS: " + strconv.Itoa(runtime.frameCounter))
		runtime.frameCounter = 0
	}
}

func (runtime *Runtime) quit() {
	runtime.window.close()
	os.Exit(0)
}

func degreesToRadians(degree float64) float64 {
	return degree * (math.Pi / 180)
}
