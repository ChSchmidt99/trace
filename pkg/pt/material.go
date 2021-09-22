package pt

import (
	"math"
	"math/rand"
)

type Material interface {
	scatter(*ray, *hit, *rand.Rand) (bool, Color)
	emittedLight() Color
}

type Light struct {
	Color Color
}

func (Light) scatter(*ray, *hit, *rand.Rand) (bool, Color) {
	return false, Color{}
}

func (l Light) emittedLight() Color {
	return l.Color
}

type Diffuse struct {
	Albedo Color
}

func (d Diffuse) scatter(ray *ray, intersec *hit, r *rand.Rand) (bool, Color) {
	scatterDirection := intersec.normal.Add(RandomUnitVector(r))

	if scatterDirection.ApproxZero() {
		scatterDirection = intersec.normal
	}

	ray.reuse(intersec.point, scatterDirection)
	return true, d.Albedo
}

func (Diffuse) emittedLight() Color {
	return NewColor(0, 0, 0)
}

type Reflective struct {
	Albedo    Color
	Diffusion float64 // diffusion in range [0,1]
}

func (d Reflective) scatter(ray *ray, intersec *hit, r *rand.Rand) (bool, Color) {
	reflected := reflect(ray.direction.Unit(), intersec.normal)
	ray.reuse(intersec.point, reflected.Add(RandomUnitVector(r).Mul(d.Diffusion)))
	return reflected.Dot(intersec.normal) > 0, d.Albedo
}

func (Reflective) emittedLight() Color {
	return NewColor(0, 0, 0)
}

type Refractive struct {
	Albedo Color
	Ratio  float64
}

func (d Refractive) scatter(ray *ray, intersec *hit, r *rand.Rand) (bool, Color) {
	refractionRatio := d.Ratio
	if intersec.frontFace {
		refractionRatio = 1 / d.Ratio
	}

	unitDir := ray.direction.Unit()
	cos_theta := math.Min(unitDir.Mul(-1).Dot(intersec.normal), 1.0)
	sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)

	cannot_refract := refractionRatio*sin_theta > 1.0

	var direction Vector3
	if cannot_refract || reflectance(cos_theta, refractionRatio) > rand.Float64() {
		direction = reflect(unitDir, intersec.normal)
	} else {
		direction = refract(unitDir, intersec.normal, refractionRatio)
	}
	ray.reuse(intersec.point, direction)
	return true, d.Albedo
}

func (Refractive) emittedLight() Color {
	return NewColor(0, 0, 0)
}

func reflect(v Vector3, n Vector3) Vector3 {
	return v.Sub(n.Mul(v.Dot(n) * 2))
}

func refract(uv Vector3, n Vector3, etai_over_etat float64) Vector3 {
	cos_theta := math.Min(uv.Mul(-1).Dot(n), 1.0)
	rOutPerp := uv.Add(n.Mul(cos_theta)).Mul(etai_over_etat)
	rOutParallel := n.Mul(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

// Schlick Approximation
func reflectance(cosine, defractionRatio float64) float64 {
	r0 := (1 - defractionRatio) / (1 + defractionRatio)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
