package pt

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var result []mortonPair

func BenchmarkSort(b *testing.B) {
	maxCode := uint64(math.Pow(float64(MORTON_SIZE), 3)) - 1
	size := 10000000
	/*
		b.Run("Go Sort", func(b *testing.B) {
			sample := generateTestSet(size, int(maxCode))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sort.Slice(sample, func(i, j int) bool {
					return sample[i].mortonCode < sample[j].mortonCode
				})
			}
		})
	*/
	b.Run("Bucket Sort", func(b *testing.B) {
		sample := generateTestSet(size, int(maxCode))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sortMortonPairs(sample, 4096, maxCode, runtime.GOMAXPROCS(0))
		}
		result = sample
	})
}

func TestBucketSort(t *testing.T) {
	size := 1000000
	maxCode := uint64(math.Pow(float64(MORTON_SIZE), 3))
	//maxCode := 1024
	sample := generateTestSet(size, int(maxCode))
	sortMortonPairs(sample, 4096, uint64(maxCode), runtime.GOMAXPROCS(0))
	for i := 1; i < len(sample); i++ {
		if sample[i].mortonCode < sample[i-1].mortonCode {
			fmt.Printf("First: %v Second: %v\n", sample[i-1], sample[i])
			assert.Fail(t, "Not sorted properly")
		}
	}
}

/*
func BenchmarkBucketSort(b *testing.B) {
	maxCode := uint64(math.Pow(float64(MORTON_SIZE), 3)) - 1
	sample := generateTestSet(1000000, int(maxCode))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sortMortonPairs(sample, 4096, maxCode, runtime.GOMAXPROCS(0))
	}
}
*/

func generateTestSet(size int, maxCode int) []mortonPair {
	out := make([]mortonPair, size)
	for i := 0; i < size; i++ {
		out[i] = mortonPair{
			primIndex:  i,
			mortonCode: uint64(rand.Intn(maxCode)),
		}
	}
	return out
}
