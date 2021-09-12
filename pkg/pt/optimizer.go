package pt

import (
	"log"
	"math"

	"changkun.de/x/bo"
)

const (
	EVAL_WIDTH  = 128
	EVAL_HEIGHT = 72
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
	x, _, err := op.o.RunSerial(func(m map[bo.Param]float64) float64 {
		alpha, delta := m[op.alphaParam], m[op.deltaParam]
		d := mapDelta(op.deltaRange, delta)
		builder.Alpha = alpha
		builder.Delta = d
		bvh := builder.BuildFromAuxilary(aux)
		cost := evalPHR(bvh, camera)
		return cost
	})
	if err != nil {
		log.Fatal(err)
	}
	return x[op.alphaParam], mapDelta(op.deltaRange, x[op.deltaParam])
}

func evalPHR(bvh BVH, camera *Camera) float64 {
	//cost := numberOfIntersections(bvh, camera)
	// TODO: revise cost function, use weights depending on scene size
	//fmt.Printf("intersections: %v cost: %v\n", cost, bvh.cost())
	//return float64(cost * bvh.size())
	return bvh.Cost()
}

func numberOfIntersections(bvh BVH, camera *Camera) int {
	sum := 0
	r := ray{}
	for y := 0; y < EVAL_HEIGHT; y++ {
		for x := 0; x < EVAL_WIDTH; x++ {
			s := float64(x) / float64(EVAL_WIDTH-1)
			t := float64(y) / float64(EVAL_HEIGHT-1)
			camera.castRayReuse(s, t, &r)
			sum += bvh.traversalSteps(r, 0.001, math.MaxFloat64)
		}
	}
	return sum
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
	for _, a := range g.Alphas {
		for _, d := range g.Deltas {
			builder.Alpha = a
			builder.Delta = d
			bvh := builder.BuildFromAuxilary(aux)
			cost := evalPHR(bvh, camera)
			//fmt.Printf("Alpha: %v Delta: %v Cost: %v\n", a, d, cost)
			if cost < minCost {
				minCost = cost
				alpha = a
				delta = d
			}
		}
	}
	return
}
