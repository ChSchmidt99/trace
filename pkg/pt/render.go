package pt

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type ClosestHitShader func(*Renderer, context, ray, *hit) Color

func DefaultClosestHitShader(renderer *Renderer, c context, r ray, h *hit) Color {
	if c.depth > renderer.maxDepth {
		return NewColor(0, 0, 0)
	}
	c.depth++
	light := h.material.emittedLight()
	// If material scatters, compute intersections with scattered ray and then call itself recursively
	if b, result := h.material.scatter(r, h, c.rand); b {
		if renderer.bvh.intersected(result.scattered, 0.0001, math.Inf(1), h) {
			return light.Add(renderer.closest(renderer, c, result.scattered, h).Blend(result.attenuation))
		} else {
			return light.Add(renderer.miss(renderer, c, result.scattered).Blend(result.attenuation))
		}
	} else {
		return light
	}
}

type MissShader func(*Renderer, context, ray) Color

func DefaultMissShader(renderer *Renderer, c context, r ray) Color {
	return NewColor(0, 0, 0)
}

type IntersectionCountShader func(count int) Color

func DefaultIntersectionCountShader(count int) Color {
	if count > 100 {
		factor := float64(count) / 150
		return NewColor(factor, 0, 0)
	}
	factor := float64(count) / 100
	return NewColor(0, 0, factor)
}

type Renderer struct {
	numCPU   int
	maxDepth int
	spp      int
	bvh      BVH
	camera   *Camera
	closest  ClosestHitShader
	miss     MissShader

	// TODO: Remove IntersectionCountShader and make different Renderer?
	intersectionCount IntersectionCountShader
}

func NewDefaultRenderer(bvh BVH, camera *Camera) *Renderer {
	return &Renderer{
		numCPU:            runtime.GOMAXPROCS(0),
		maxDepth:          2,
		bvh:               bvh,
		spp:               3000,
		camera:            camera,
		closest:           DefaultClosestHitShader,
		miss:              DefaultMissShader,
		intersectionCount: DefaultIntersectionCountShader,
	}
}

type context struct {
	rand  *rand.Rand
	depth int
}

func (r *Renderer) IntersectionHeatMapToBuffer(buff Buffer) {
	jobs := make(chan int, buff.h())
	wg := sync.WaitGroup{}
	wg.Add(r.numCPU)
	width := buff.w()
	height := buff.h()
	for i := 0; i < r.numCPU; i++ {
		go func(c context, w, h int) {
			ray := ray{
				// ATTENTION! When camera is moved, this origin needs to be changed too!
				origin: r.camera.orientation.origin,
			}
			for y := range jobs {
				for x := 0; x < w; x++ {
					u := (float64(x) + c.rand.Float64()) / float64(w-1)
					v := (float64(y) + c.rand.Float64()) / float64(h-1)
					r.camera.castRayReuse(u, v, &ray)
					count := r.bvh.intersectionTests(ray, 0.001, math.Inf(1))
					buff.addSample(x, y, r.intersectionCount(count))
				}
			}
			wg.Done()
		}(context{
			rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		}, width, height)
	}
	// TODO: Check if doing spp per pixel instead of per image is better
	for i := 0; i < r.spp; i++ {
		for y := 0; y < height; y++ {
			jobs <- y
		}
	}
	close(jobs)
	wg.Wait()
}

func (r *Renderer) RenderToBuffer(buff Buffer) {
	jobs := make(chan int, buff.h())
	wg := sync.WaitGroup{}
	wg.Add(r.numCPU)
	width := buff.w()
	height := buff.h()
	for i := 0; i < r.numCPU; i++ {
		go func(c context, w, h int) {
			ray := ray{
				// ATTENTION! When camera is moved, this origin needs to be changed too!
				origin: r.camera.orientation.origin,
			}
			hit := hit{}
			for y := range jobs {
				for x := 0; x < w; x++ {
					u := (float64(x) + c.rand.Float64()) / float64(w-1)
					v := (float64(y) + c.rand.Float64()) / float64(h-1)
					r.camera.castRayReuse(u, v, &ray)
					if r.bvh.intersected(ray, 0.001, math.Inf(1), &hit) {
						buff.addSample(x, y, r.closest(r, c, ray, &hit))
					} else {
						buff.addSample(x, y, r.miss(r, c, ray))
					}
				}
			}
			wg.Done()
		}(context{
			rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		}, width, height)
	}
	// TODO: Check if doing spp per pixel instead of per image is better
	for i := 0; i < r.spp; i++ {
		for y := 0; y < height; y++ {
			jobs <- y
		}
	}
	close(jobs)
	wg.Wait()
}
