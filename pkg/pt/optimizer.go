package pt

import "math"

// TODO: Reevaluate these numbers
const (
	PRIM_THRESH = 300000
	RES_THRESH  = 100000
)

type Optimizer interface {
	OptimizedPHRparams(aux BVH, branching int, threads int) (alpha float64, delta float64)
}

type gridOptimizer struct {
	Alphas []float64
	Deltas []float64
	pixels int
}

func NewDefaultGridOptimizer(frameWidth int, frameHeight int) gridOptimizer {
	return NewGridOptimizer([]float64{0.4, 0.45, 0.5}, []float64{6, 7, 8, 9}, frameWidth, frameHeight)
}

func NewGridOptimizer(alphas []float64, deltas []float64, frameWidth int, frameHeight int) gridOptimizer {
	return gridOptimizer{
		Alphas: alphas,
		Deltas: deltas,
		pixels: frameWidth * frameHeight,
	}
}

type evaluation struct {
	sahCost   float64
	buildCost int
	alpha     float64
	delta     float64
}

func (g gridOptimizer) OptimizedPHRparams(aux BVH, branching int, threads int) (alpha float64, delta float64) {
	builder := NewPHRBuilder(aux.prims, 0, 0, branching, threads)
	evaluations := make([]evaluation, 0, len(g.Alphas)*len(g.Deltas))
	// Run all alpha-delta combinations
	for _, a := range g.Alphas {
		for _, d := range g.Deltas {
			builder.Alpha = a
			builder.Delta = d
			bvh, buildCost := builder.BuildWithCost(aux)
			evaluations = append(evaluations, evaluation{
				sahCost:   bvh.Cost(),
				buildCost: buildCost,
				alpha:     a,
				delta:     d,
			})
		}
	}

	// Determine max values
	maxBuildCost := 0
	maxSAHCost := 0.0
	for _, eval := range evaluations {
		if eval.buildCost > maxBuildCost {
			maxBuildCost = eval.buildCost
		}
		if eval.sahCost > maxSAHCost {
			maxSAHCost = eval.sahCost
		}
	}

	// Evaluate results
	minCost := math.MaxFloat64
	var bestEvaluation evaluation
	o := omega(len(aux.prims), g.pixels)
	for _, eval := range evaluations {
		cost := evalPHR(eval.sahCost, maxSAHCost, eval.buildCost, maxBuildCost, o)
		if cost < float64(minCost) {
			minCost = cost
			bestEvaluation = eval
		}
	}
	return bestEvaluation.alpha, bestEvaluation.delta
}

func omega(primitives int, pixels int) float64 {
	return (2.0*(float64(primitives)/PRIM_THRESH) + float64(pixels)/RES_THRESH) / 3
}

func evalPHR(SAHcost, maxSAHcost float64, buildCost int, maxBuildCost int, omega float64) float64 {
	return (float64(buildCost) / float64(maxBuildCost)) + omega*(SAHcost/maxSAHcost)
}

/*

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

func (op bayesianOptimizer) OptimizedPHRparams(aux BVH, camera *Camera, branching int, threads int) (alpha float64, delta float64) {
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

func mapDelta(deltaRange [2]int, t float64) int {
	abs := deltaRange[1] - deltaRange[0]
	return deltaRange[0] + int(float64(abs)*t)
}
*/
