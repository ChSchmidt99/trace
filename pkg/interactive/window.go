package raytracer

import (
	"github/chschmidt99/pt/pkg/pt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	buffer *pt.FrameBuffer
	name   string
	window *sdl.Window
}

type EventType int

const (
	DOWN EventType = 0
	UP   EventType = 1
)

type KeyEvent struct {
	keyCode   sdl.Keycode
	eventType EventType
	timestamp time.Time
}

type MouseEvent struct {
	x int // relative x movement
	y int // relative y movement
}

// InitWindow() needs to be called to actually instantiate the window
func NewWindow(height int, aspectRatio float64, name string) *Window {
	return &Window{
		buffer: pt.NewFrameBufferAR(height, aspectRatio),
		name:   name,
	}
}

// initializes SDL and opens up a new Window
func (win *Window) InitWindow() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow(win.name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(win.buffer.Width()), int32(win.buffer.Height()), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	sdl.SetRelativeMouseMode(true)

	win.window = window
}

func (win *Window) Close() {
	win.window.Destroy()
	sdl.Quit()
}

// Starts the main loop. Each iteration the updateHandler will be called,
// then all inputs are checked and possible events passed to the eventHandler.
func (win *Window) Run(updateHandler func(window *Window), eventHandler func(event interface{})) {
	running := true
	for running {
		updateHandler(win)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch value := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if value.GetType() == sdl.KEYDOWN {
					if value.Repeat == 0 {
						eventHandler(&KeyEvent{
							keyCode:   value.Keysym.Sym,
							eventType: DOWN,
							timestamp: time.Unix(int64(value.Timestamp), 0),
						})
					}
				} else {
					eventHandler(&KeyEvent{
						keyCode:   value.Keysym.Sym,
						eventType: UP,
						timestamp: time.Unix(int64(value.Timestamp), 0),
					})
				}
			case *sdl.MouseMotionEvent:
				eventHandler(&MouseEvent{
					x: int(value.XRel),
					y: int(value.YRel),
				})
			}
		}
	}
}

// Write new values to color buffer and then update the view
func (win *Window) Refresh() {
	surface, err := win.window.GetSurface()
	if err != nil {
		panic(err)
	}
	width := win.buffer.Width()
	height := win.buffer.Height()
	for i := 0; i < width*height; i++ {
		x := i % width
		y := (height - (i / width)) - 1
		color := win.buffer.GoColor(i)
		surface.Set(x, y, color)
	}
	win.window.UpdateSurface()
}

func (win *Window) close() {
	win.window.Destroy()
}
