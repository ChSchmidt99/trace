package pt

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type ClosestHitShader func(*Ray, *Intersection) Color

type MissShader func(*Ray) Color

type Buffer struct {
	Width  int
	Height int
	Buff   []Color
}

func NewBufferAspect(height int, aspect float64) *Buffer {
	width := float64(height) * aspect
	return NewBuffer(int(width), height)
}

func NewBuffer(width, height int) *Buffer {
	return &Buffer{
		Width:  width,
		Height: height,
		Buff:   make([]Color, width*height),
	}
}

func (b *Buffer) setPixel(x, y int, c Color) {
	b.Buff[y*b.Width+x] = c
}

type Renderer struct {
	numCPU   int
	maxDepth int
	bvh      BVH
	camera   *Camera
	closest  ClosestHitShader
	miss     MissShader
}

func NewDefaultRenderer(bvh BVH, camera *Camera) *Renderer {
	return &Renderer{
		numCPU:   runtime.GOMAXPROCS(0),
		maxDepth: 5,
		bvh:      bvh,
		camera:   camera,
		closest: func(r *Ray, i *Intersection) Color {
			return NewColor(1, 0, 0)
		},
		miss: func(r *Ray) Color {
			return NewColor(1, 1, 1)
		},
	}
}

type context struct {
	rand *rand.Rand
}

// Render to a buffer that is already allocated in the correct size
func (r *Renderer) RenderToBuffer(buff *Buffer) {

	jobs := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(r.numCPU)
	for i := 0; i < r.numCPU; i++ {
		go func(c context) {
			for line := range jobs {
				r.renderLine(&c, line, buff)
			}
			wg.Done()
		}(context{
			rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		})
	}

	for y := 0; y < buff.Height; y++ {
		jobs <- y
	}

	close(jobs)
	wg.Wait()
}

func (r *Renderer) renderLine(c *context, y int, buffer *Buffer) {
	for x := 0; x < buffer.Width; x++ {
		u := (float64(x) + c.rand.Float64()) / float64(buffer.Width-1)
		v := (float64(y) + c.rand.Float64()) / float64(buffer.Height-1)
		ray := r.camera.castRay(u, v)
		if intersection := r.bvh.intersected(ray, 0.001, math.Inf(1)); intersection != nil {
			buffer.setPixel(x, y, r.closest(ray, intersection))
		} else {
			buffer.setPixel(x, y, r.miss(ray))
		}
	}
}
