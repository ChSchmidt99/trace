package pt

import (
	"log"
	"math"

	"changkun.de/x/bo"
)

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

func NewDefaultGridOptimizer(frameWidth, frameHeight int) gridOptimizer {
	return NewGridOptimizer([]float64{0.4, 0.45, 0.5}, []float64{6, 7, 8, 9}, frameWidth, frameHeight)
}

func NewGridOptimizer(alphas []float64, deltas []float64, frameWidth, frameHeight int) gridOptimizer {
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

type bayesianOptimizer struct {
	alphaParam bo.UniformParam
	deltaParam bo.UniformParam
	alphaRange [2]float64
	deltaRange [2]float64
	o          *bo.Optimizer
	pixels     int
}

func NewDefaultBayesianOptimizer(frameWidth, frameHeight int) bayesianOptimizer {
	alphaRange := [2]float64{0.4, 0.5}
	deltaRange := [2]float64{6, 9}
	return NewBayesianOptimizer(alphaRange, deltaRange, frameWidth, frameHeight)
}

func NewBayesianOptimizer(alphaRange [2]float64, deltaRange [2]float64, frameWidth, frameHeight int) bayesianOptimizer {
	alpha := bo.UniformParam{
		Min: alphaRange[0],
		Max: alphaRange[1],
	}
	delta := bo.UniformParam{
		Min: deltaRange[0],
		Max: deltaRange[1],
	}
	o := bo.NewOptimizer([]bo.Param{alpha, delta})
	return bayesianOptimizer{
		alphaParam: alpha,
		deltaParam: delta,
		o:          o,
		pixels:     frameWidth * frameHeight,
	}
}

func (op bayesianOptimizer) OptimizedPHRparams(aux BVH, branching int, threads int) (alpha float64, delta float64) {

	builder := NewPHRBuilder(aux.prims, op.alphaRange[0], op.deltaRange[0], branching, threads)
	bvh := builder.BuildFromAuxilary(aux)
	maxSAHcost := bvh.Cost()

	builder.Alpha = op.alphaRange[1]
	builder.Delta = op.deltaRange[1]
	_, maxBuildCost := builder.BuildWithCost(aux)

	omega := omega(len(aux.prims), op.pixels)
	x, _, err := op.o.RunSerial(func(m map[bo.Param]float64) float64 {
		alpha, delta := m[op.alphaParam], m[op.deltaParam]
		builder.Alpha = alpha
		builder.Delta = delta
		bvh, buildCost := builder.BuildWithCost(aux)
		cost := evalPHR(bvh.Cost(), maxSAHcost, buildCost, maxBuildCost, omega)
		return cost
	})
	if err != nil {
		log.Fatal(err)
	}
	return x[op.alphaParam], x[op.deltaParam]
}
