package pt

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type ClosestHitShader func(*Renderer, context, *ray, *hit) Color

func DefaultClosestHitShader(renderer *Renderer, c context, r *ray, h *hit) Color {
	if c.depth > renderer.maxDepth {
		return NewColor(1, 1, 1)
	}
	c.depth++
	// If material scatters, compute intersections with scattered ray and then call itself recursively
	if b, result := h.material.scatter(r, h, c.rand); b {
		hit := renderer.bvh.intersected(result.scattered, 0.001, math.Inf(1))
		if hit == nil {
			return renderer.miss(renderer, c, result.scattered).Blend(result.attenuation)
		}
		return renderer.closest(renderer, c, result.scattered, hit).Blend(result.attenuation)
	} else {
		return result.attenuation
	}
}

type MissShader func(*Renderer, context, *ray) Color

func DefaultMissShader(renderer *Renderer, c context, r *ray) Color {
	return NewColor(1, 1, 1)
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
		closest:  DefaultClosestHitShader,
		miss:     DefaultMissShader,
	}
}

// TODO: Check if it's faster to use renderContex ptr insted of value
type context struct {
	rand  *rand.Rand
	depth int
}

// Render to a buffer that is already allocated in the correct size
func (r *Renderer) RenderToBuffer(buff *Buffer) {
	jobs := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(r.numCPU)
	for i := 0; i < r.numCPU; i++ {
		go func(c context) {
			for y := range jobs {
				r.renderLine(c, y, buff)
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

func (r *Renderer) renderLine(c context, y int, buffer *Buffer) {
	for x := 0; x < buffer.Width; x++ {
		u := (float64(x) + c.rand.Float64()) / float64(buffer.Width-1)
		v := (float64(y) + c.rand.Float64()) / float64(buffer.Height-1)
		ray := r.camera.castray(u, v)
		if intersection := r.bvh.intersected(ray, 0.001, math.Inf(1)); intersection != nil {
			buffer.setPixel(x, y, r.closest(r, c, ray, intersection))
		} else {
			buffer.setPixel(x, y, r.miss(r, c, ray))
		}
	}
}
