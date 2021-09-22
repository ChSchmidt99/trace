package pt

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Renderer interface {
	RenderToBuffer(buff Buffer)
}

type ImageRenderer struct {
	NumCPU   int
	MaxDepth int
	Spp      int
	Bvh      BVH
	Camera   *Camera
	Closest  ClosestHitShader
	Miss     MissShader
	Sampling Sampling
	Verbose  bool
}

func NewDefaultRenderer(bvh BVH, camera *Camera) *ImageRenderer {
	return &ImageRenderer{
		NumCPU:   runtime.GOMAXPROCS(0),
		MaxDepth: 2,
		Bvh:      bvh,
		Spp:      300,
		Camera:   camera,
		Closest:  DefaultClosestHitShader,
		Miss:     DefaultMissShader,
		Sampling: RandomSampling,
		Verbose:  false,
	}
}

func NewNoLightRenderer(bvh BVH, camera *Camera) *ImageRenderer {
	return &ImageRenderer{
		NumCPU:   runtime.GOMAXPROCS(0),
		MaxDepth: 2,
		Bvh:      bvh,
		Spp:      300,
		Camera:   camera,
		Closest:  UnlitClosestHitShader,
		Miss:     WhiteMissShader,
		Sampling: RandomSampling,
		Verbose:  false,
	}
}

func NewBenchmarkRenderer(bvh BVH, camera *Camera) *ImageRenderer {
	return &ImageRenderer{
		NumCPU:   runtime.GOMAXPROCS(0),
		MaxDepth: 5,
		Bvh:      bvh,
		Spp:      1,
		Camera:   camera,
		Closest:  UnlitClosestHitShader,
		Miss:     WhiteMissShader,
		Sampling: RandomSampling,
		Verbose:  false,
	}
}

type context struct {
	rand  *rand.Rand
	depth int
}

func (r *ImageRenderer) RenderToBuffer(buff Buffer) {
	r.log("Started rendering\n")
	jobs := make(chan int, buff.h())
	wg := sync.WaitGroup{}
	wg.Add(r.NumCPU)
	width := buff.w()
	height := buff.h()
	for i := 0; i < r.NumCPU; i++ {
		go func(c context, w, h int) {
			ray := ray{
				origin: r.Camera.orientation.origin,
			}
			hit := hit{}
			for y := range jobs {
				for x := 0; x < w; x++ {
					u, v := r.Sampling(c, x, y, w, h)
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
	for i := 0; i < r.Spp; i++ {
		for y := 0; y < height; y++ {
			jobs <- y
		}
		r.log("Finished pass %v\n", i)
	}
	close(jobs)
	wg.Wait()
}

// TODO: Add Render function for image (RenderIncremental) that saves to image in interval and Renderer that finishes each pixel with all spp at once

func (r *ImageRenderer) log(message string, a ...interface{}) {
	if r.Verbose {
		fmt.Printf(message, a...)
	}
}

type HeatMapRenderer struct {
	NumCPU            int
	Bvh               BVH
	Camera            *Camera
	Threshold         int
	intersectionCount TraversalCountShader
}

func NewHeatMapRenderer(bvh BVH, camera *Camera, threshold int) *HeatMapRenderer {
	return &HeatMapRenderer{
		NumCPU:            runtime.GOMAXPROCS(0),
		Bvh:               bvh,
		Camera:            camera,
		Threshold:         threshold,
		intersectionCount: DefaultTraversalCountShader,
	}
}

func (r *HeatMapRenderer) RenderToBuffer(buff Buffer) {
	jobs := make(chan int, buff.h())
	wg := sync.WaitGroup{}
	wg.Add(r.NumCPU)
	width := buff.w()
	height := buff.h()
	for i := 0; i < r.NumCPU; i++ {
		go func(c context, w, h int) {
			ray := ray{
				origin: r.Camera.orientation.origin,
			}
			for y := range jobs {
				for x := 0; x < w; x++ {
					u := float64(x) / float64(w-1)
					v := float64(y) / float64(h-1)
					r.Camera.castRayReuse(u, v, &ray)
					count := r.Bvh.traversalSteps(ray, 0.001, math.MaxFloat64)
					buff.addSample(x, y, r.intersectionCount(r, count))
				}
			}
			wg.Done()
		}(context{
			rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		}, width, height)
	}
	for y := 0; y < height; y++ {
		jobs <- y
	}
	close(jobs)
	wg.Wait()
}

type Sampling func(c context, x, y, w, h int) (u, v float64)

func RandomSampling(c context, x, y, w, h int) (u, v float64) {
	u = (float64(x) + c.rand.Float64()) / float64(w-1)
	v = (float64(y) + c.rand.Float64()) / float64(h-1)
	return
}

type ClosestHitShader func(*ImageRenderer, context, ray, *hit) Color

func DefaultClosestHitShader(renderer *ImageRenderer, c context, r ray, h *hit) Color {
	if c.depth > renderer.MaxDepth {
		return NewColor(0, 0, 0)
	}
	c.depth++
	light := h.material.emittedLight()
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

// TODO: Needed?
func UnlitClosestHitShader(renderer *ImageRenderer, c context, r ray, h *hit) Color {
	if c.depth > renderer.MaxDepth {
		return NewColor(0, 0, 0)
	}
	c.depth++
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

type MissShader func(*ImageRenderer, context, ray) Color

func DefaultMissShader(renderer *ImageRenderer, c context, r ray) Color {
	return NewColor(0, 0, 0)
}

func WhiteMissShader(renderer *ImageRenderer, c context, r ray) Color {
	return NewColor(1, 1, 1)
}

func SkyMissShader(renderer *ImageRenderer, c context, r ray) Color {
	unit := r.direction.Unit()
	t := 0.5 * (unit.Y + 1)
	white := NewColor(0.8, 0.8, 0.8)
	blue := NewColor(0.25, 0.35, 0.5)
	return white.Scale(1.0 - t).Add(blue.Scale(t))
}

type TraversalCountShader func(renderer *HeatMapRenderer, count int) Color

func DefaultTraversalCountShader(renderer *HeatMapRenderer, count int) Color {
	if count > renderer.Threshold {
		factor := float64(count) / float64(renderer.Threshold) * 2
		return NewColor(factor, 0, 0)
	}
	factor := float64(count) / (float64(renderer.Threshold) * 1.25)
	return NewColor(0, factor, 1-factor)

}
