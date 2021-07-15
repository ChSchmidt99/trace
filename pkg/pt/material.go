package pt

import (
	"math"
	"math/rand"
)

type scatterResult struct {
	scattered   ray
	attenuation Color
}

type Material interface {
	scatter(ray, *hit, *rand.Rand) (bool, scatterResult)
}

type Diffuse struct {
	Albedo Color
}

func (d *Diffuse) scatter(ray ray, intersec *hit, r *rand.Rand) (bool, scatterResult) {
	// TODO: Revise random ray generation
	scatterDirection := intersec.normal.Add(RandomUnitVector(r))
	if scatterDirection.ApproxZero() {
		scatterDirection = intersec.normal
	}
	ray.reuse(intersec.point, scatterDirection)
	return true, scatterResult{
		scattered:   ray,
		attenuation: d.Albedo,
	}
}

type Reflective struct {
	Albedo    Color
	Diffusion float64 // diffusion in range [0,1]
}

func (d *Reflective) scatter(ray ray, intersec *hit, r *rand.Rand) (bool, scatterResult) {
	reflected := reflect(ray.direction.Unit(), intersec.normal)
	ray.reuse(intersec.point, reflected.Add(RandomUnitVector(r).Mul(d.Diffusion)))
	return reflected.Dot(intersec.normal) > 0, scatterResult{
		scattered:   ray,
		attenuation: d.Albedo,
	}
}

type Refractive struct {
	Albedo Color
	Ratio  float64
}

func (d *Refractive) scatter(ray ray, intersec *hit, r *rand.Rand) (bool, scatterResult) {
	defractionRatio := d.Ratio
	if intersec.frontFace {
		defractionRatio = 1 / d.Ratio
	}

	unitDir := ray.direction.Unit()
	cos_theta := math.Min(unitDir.Mul(-1).Dot(intersec.normal), 1.0)
	sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)

	cannot_refract := defractionRatio*sin_theta > 1.0

	var direction Vector3
	if cannot_refract || reflectance(cos_theta, defractionRatio) > rand.Float64() {
		direction = reflect(unitDir, intersec.normal)
	} else {
		direction = refract(unitDir, intersec.normal, defractionRatio)
	}
	ray.reuse(intersec.point, direction)
	return true, scatterResult{
		scattered:   ray,
		attenuation: d.Albedo,
	}
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
