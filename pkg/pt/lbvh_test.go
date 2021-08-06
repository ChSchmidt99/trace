package pt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPrimitive struct {
	center Vector3
	size   float64
}

func (*testPrimitive) intersected(ray ray, tMin, tMax float64, hitOut *hit) bool {
	return false
}

func (p *testPrimitive) bounding() aabb {
	min := p.center.Sub(NewVector3(p.size/2.0, p.size/2.0, p.size/2.0))
	max := p.center.Add(NewVector3(p.size/2.0, p.size/2.0, p.size/2.0))
	return newAABB(min, max)
}

func (p *testPrimitive) transformed(t Matrix4) Primitive {
	return p
}

func TestAssignMortonCodes(t *testing.T) {
	enclosing := newAABB(NewVector3(-2, -2, -2), NewVector3(2, 2, 2))
	prim1 := testPrimitive{
		center: NewVector3(-1.75, -1.75, -1.75),
	}
	prim2 := testPrimitive{
		center: NewVector3(1.75, 1.75, 1.75),
	}
	prim3 := testPrimitive{
		center: NewVector3(0, 0, 0),
	}
	expected := []mortonPair{
		{
			primIndex:  0,
			mortonCode: 0,
		},
		{
			primIndex:  1,
			mortonCode: 56,
		},
		{
			primIndex:  2,
			mortonCode: 7,
		},
	}
	pairs := assignMortonCodes([]Primitive{&prim1, &prim2, &prim3}, enclosing, 4, 4)
	assert.Equal(t, expected, pairs)
}
