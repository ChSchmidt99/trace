package pt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransposed(t *testing.T) {
	expected := Matrix4{
		values: [4][4]float64{
			{1, 5, 9, 13},
			{2, 6, 10, 14},
			{3, 7, 11, 15},
			{4, 8, 12, 16},
		},
	}
	m := Matrix4{
		values: [4][4]float64{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 16},
		},
	}
	assert.Equal(t, expected, m.Transpose())
}

func TestInverse(t *testing.T) {
	expected := Matrix4{
		values: [4][4]float64{
			{0.25, 0.25, 5, -3.25},
			{0.75, -0.25, 0, 0.25},
			{-1.5, 0.5, 3, -1.5},
			{0, 0, -2, 1},
		},
	}
	m := Matrix4{
		values: [4][4]float64{
			{1, 3, 1, 4},
			{3, 9, 5, 15},
			{0, 2, 1, 1},
			{0, 4, 2, 3},
		},
	}
	assert.Equal(t, expected, m.Inverse())
}
