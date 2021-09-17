package pt

import (
	"fmt"
	"log"
	"math"
	"time"

	"changkun.de/x/bo"
)

const (
	EVAL_WIDTH  = 256
	EVAL_HEIGHT = 144
	//EVAL_WIDTH  = 256
	//EVAL_HEIGHT = 144
)

type Optimizer interface {
	OptimizedPHRparams(aux BVH, camera *Camera, branching int, threads int) (alpha float64, delta int)
}

type bayesianOptimizer struct {
	alphaParam bo.UniformParam
	deltaParam bo.UniformParam
	deltaRange [2]int
	o          *bo.Optimizer
}

func NewBayesianOptimizer(alphaRange [2]float64, deltaRange [2]int) bayesianOptimizer {
	alpha := bo.UniformParam{
		Min: alphaRange[0],
		Max: alphaRange[1],
	}
	delta := bo.UniformParam{
		Min: 0,
		Max: 1,
	}
	o := bo.NewOptimizer([]bo.Param{alpha, delta})
	return bayesianOptimizer{
		alphaParam: alpha,
		deltaParam: delta,
		deltaRange: deltaRange,
		o:          o,
	}
}

// TODO: Also optimize Branching factor?
func (op bayesianOptimizer) OptimizedPHRparams(aux BVH, camera *Camera, branching int, threads int) (alpha float64, delta int) {
	builder := NewPHRBuilder(aux.prims, 0, 0, branching, threads)
	e := evaluater{
		set: false,
	}
	x, _, err := op.o.RunSerial(func(m map[bo.Param]float64) float64 {
		alpha, delta := m[op.alphaParam], m[op.deltaParam]
		d := mapDelta(op.deltaRange, delta)
		builder.Alpha = alpha
		builder.Delta = d
		start := time.Now()
		bvh := builder.BuildFromAuxilary(aux)
		cost := e.evalPHR(bvh, camera, time.Since(start))
		return cost
	})
	if err != nil {
		log.Fatal(err)
	}
	return x[op.alphaParam], mapDelta(op.deltaRange, x[op.deltaParam])
}

// Stores the initial evaluation to use as a base line
// TODO: Add weight depending on canvas size and scene complexity
// TODO: Better metric than build time (read build complexity from params)
type evaluater struct {
	cost      float64
	buildTime time.Duration
	set       bool
}

func (e *evaluater) evalPHR(bvh BVH, camera *Camera, buildTime time.Duration) float64 {
	if !e.set {
		e.cost = AvgRayCost(bvh, camera)
		e.buildTime = buildTime
		e.set = true
		return 2
	}
	cost := AvgRayCost(bvh, camera)
	eval := (cost / e.cost) + (float64(buildTime) / float64(e.buildTime))
	fmt.Printf("Build Time: %v Cost: %v Evaluation: %v\n", buildTime, cost, eval)
	return eval
}

// TODO: Make private
func AvgRayCost(bvh BVH, camera *Camera) float64 {
	costSum := 0.0
	r := ray{}
	for y := 0; y < EVAL_HEIGHT; y++ {
		for x := 0; x < EVAL_WIDTH; x++ {
			s := float64(x) / float64(EVAL_WIDTH-1)
			t := float64(y) / float64(EVAL_HEIGHT-1)
			camera.castRayReuse(s, t, &r)
			costSum += bvh.rayCost(r, 0.001, math.MaxFloat64)
		}
	}
	return costSum
}

func mapDelta(deltaRange [2]int, t float64) int {
	abs := deltaRange[1] - deltaRange[0]
	return deltaRange[0] + int(float64(abs)*t)
}

type gridOptimizer struct {
	Alphas []float64
	Deltas []int
}

func NewGridOptimizer(alphas []float64, deltas []int) gridOptimizer {
	return gridOptimizer{
		Alphas: alphas,
		Deltas: deltas,
	}
}

func (g gridOptimizer) OptimizedPHRparams(aux BVH, camera *Camera, branching int, threads int) (alpha float64, delta int) {
	builder := NewPHRBuilder(aux.prims, 0, 0, branching, threads)
	minCost := math.MaxFloat64
	e := evaluater{
		set: false,
	}
	for _, a := range g.Alphas {
		for _, d := range g.Deltas {
			builder.Alpha = a
			builder.Delta = d
			start := time.Now()
			bvh := builder.BuildFromAuxilary(aux)
			cost := e.evalPHR(bvh, camera, time.Since(start))
			if cost < minCost {
				minCost = cost
				alpha = a
				delta = d
			}
		}
	}
	return
}
