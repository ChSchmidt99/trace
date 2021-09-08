package pt

import (
	"math"
	"sort"
	"sync"
)

type bucket []mortonPair

type mergeJob struct {
	index   int
	buckets []bucket
}

// Parallel Bucket sort
func sortMortonPairs(pairs []mortonPair, numberOfBuckets int, maxMorton uint64, threads int) []mortonPair {
	// Put pairs into buckets
	batchSize := int(math.Ceil(float64(len(pairs)) / float64(threads)))

	// TODO: Test against shared bucket between threads
	bucketCollection := make([][]bucket, 0, threads)
	wg := sync.WaitGroup{}
	for i := 0; i < threads; i++ {
		start := i * batchSize
		if start >= len(pairs) {
			break
		}
		end := int(math.Min(float64(start+batchSize), float64(len(pairs))))
		bucketCollection = append(bucketCollection, make([]bucket, numberOfBuckets))
		wg.Add(1)
		go func(input []mortonPair, threadNumber int) {
			for _, pair := range input {
				index := (uint64(numberOfBuckets) * pair.mortonCode) / maxMorton
				bucketCollection[threadNumber][index] = append(bucketCollection[threadNumber][index], pair)
			}
			wg.Done()
		}(pairs[start:end], i)
	}
	wg.Wait()
	sorted := make([]bucket, numberOfBuckets)
	jobs := make(chan mergeJob, threads)
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			for job := range jobs {
				sorted[job.index] = mergeBuckets(job.buckets)
			}
			wg.Done()
		}()
	}

	for i := 0; i < numberOfBuckets; i++ {
		job := mergeJob{
			index: i,
		}
		for _, buck := range bucketCollection {
			job.buckets = append(job.buckets, buck[i])
		}
		jobs <- job
	}
	close(jobs)
	wg.Wait()

	// TODO: Parallelize
	index := 0
	for _, buck := range sorted {
		for _, pair := range buck {
			pairs[index] = pair
			index++
		}
	}
	return pairs
}

// Merges n buckets with the same bucket index
func mergeBuckets(buckets []bucket) bucket {
	out := make([]mortonPair, 0)
	for _, bucket := range buckets {
		out = append(out, bucket...)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].mortonCode < out[j].mortonCode
	})
	return out
}
