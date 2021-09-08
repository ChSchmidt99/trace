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
	if c.depth > renderer.MaxDepth {
		return NewColor(0, 0, 0)
	}
	c.depth++
	light := h.material.emittedLight()
	// If material scatters, compute intersections with scattered ray and then call itself recursively
	if b, attenuation := h.material.scatter(&r, h, c.rand); b {
		if renderer.Bvh.intersected(r, 0.0001, math.Inf(1), h) {
			return light.Add(renderer.Closest(renderer, c, r, h).Blend(attenuation))
		} else {
			return light.Add(renderer.Miss(renderer, c, r).Blend(attenuation))
		}
	} else {
		return light
	}
}

func UnlitClosestHitShader(renderer *Renderer, c context, r ray, h *hit) Color {
	if c.depth > renderer.MaxDepth {
		return NewColor(0, 0, 0)
	}
	c.depth++
	// If material scatters, compute intersections with scattered ray and then call itself recursively
	if b, attenuation := h.material.scatter(&r, h, c.rand); b {
		if renderer.Bvh.intersected(r, 0.0001, math.Inf(1), h) {
			return renderer.Closest(renderer, c, r, h).Blend(attenuation)
		} else {
			return renderer.Miss(renderer, c, r).Blend(attenuation)
		}
	} else {
		return attenuation
	}
}

type MissShader func(*Renderer, context, ray) Color

func DefaultMissShader(renderer *Renderer, c context, r ray) Color {
	return NewColor(0, 0, 0)
}

func WhiteMissShader(renderer *Renderer, c context, r ray) Color {
	return NewColor(1, 1, 1)
}

func SkyMissShader(renderer *Renderer, c context, r ray) Color {
	unit := r.direction.Unit()
	t := 0.5 * (unit.Y + 1)
	white := NewColor(0.8, 0.8, 0.8)
	blue := NewColor(0.25, 0.35, 0.5)
	return white.Scale(1.0 - t).Add(blue.Scale(t))
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
	NumCPU   int
	MaxDepth int
	Spp      int
	Bvh      BVH
	Camera   *Camera
	Closest  ClosestHitShader
	Miss     MissShader

	// TODO: Add verbose mode and incremental saving to image
	// TODO: Remove IntersectionCountShader and make different Renderer?
	intersectionCount IntersectionCountShader
}

func NewDefaultRenderer(bvh BVH, camera *Camera) *Renderer {
	return &Renderer{
		NumCPU:            runtime.GOMAXPROCS(0),
		MaxDepth:          2,
		Bvh:               bvh,
		Spp:               300,
		Camera:            camera,
		Closest:           DefaultClosestHitShader,
		Miss:              DefaultMissShader,
		intersectionCount: DefaultIntersectionCountShader,
	}
}

func NewNoLightRenderer(bvh BVH, camera *Camera) *Renderer {
	return &Renderer{
		NumCPU:            runtime.GOMAXPROCS(0),
		MaxDepth:          2,
		Bvh:               bvh,
		Spp:               300,
		Camera:            camera,
		Closest:           UnlitClosestHitShader,
		Miss:              WhiteMissShader,
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
	wg.Add(r.NumCPU)
	width := buff.w()
	height := buff.h()
	for i := 0; i < r.NumCPU; i++ {
		go func(c context, w, h int) {
			ray := ray{
				// ATTENTION! When camera is moved, this origin needs to be changed too!
				origin: r.Camera.orientation.origin,
			}
			for y := range jobs {
				for x := 0; x < w; x++ {
					u := (float64(x) + c.rand.Float64()) / float64(w-1)
					v := (float64(y) + c.rand.Float64()) / float64(h-1)
					r.Camera.castRayReuse(u, v, &ray)
					count := r.Bvh.intersectionTests(ray, 0.001, math.Inf(1))
					buff.addSample(x, y, r.intersectionCount(count))
				}
			}
			wg.Done()
		}(context{
			rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		}, width, height)
	}
	// TODO: Check if doing spp per pixel instead of per image is better
	for i := 0; i < r.Spp; i++ {
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
	wg.Add(r.NumCPU)
	width := buff.w()
	height := buff.h()
	for i := 0; i < r.NumCPU; i++ {
		go func(c context, w, h int) {
			ray := ray{
				// ATTENTION! When camera is moved, this origin needs to be changed too!
				origin: r.Camera.orientation.origin,
			}
			hit := hit{}
			for y := range jobs {
				for x := 0; x < w; x++ {
					u := (float64(x) + c.rand.Float64()) / float64(w-1)
					v := (float64(y) + c.rand.Float64()) / float64(h-1)
					r.Camera.castRayReuse(u, v, &ray)
					if r.Bvh.intersected(ray, 0.001, math.Inf(1), &hit) {
						buff.addSample(x, y, r.Closest(r, c, ray, &hit))
					} else {
						buff.addSample(x, y, r.Miss(r, c, ray))
					}
				}
			}
			wg.Done()
		}(context{
			rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		}, width, height)
	}
	// TODO: Check if doing spp per pixel instead of per image is better
	for i := 0; i < r.Spp; i++ {
		for y := 0; y < height; y++ {
			// TODO: Check if not using worker pattern and instead rendering calculated lines is faster
			jobs <- y
		}
	}
	close(jobs)
	wg.Wait()
}
