package pt

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

type bucket []mortonPair

type mergeJob struct {
	index   int
	buckets []bucket
	out     []mortonPair
}

// Parallel Bucket sort
func sortMortonPairs(pairs []mortonPair, numberOfBuckets int, maxMorton uint64, threads int) []mortonPair {
	// Put pairs into buckets
	batchSize := int(math.Ceil(float64(len(pairs)) / float64(threads)))
	bucketCollection := make([][]bucket, 0, threads)

	// Stores how many pairs are stored in the corresponding bucket
	bucketEntries := make([]int32, numberOfBuckets)

	wg := sync.WaitGroup{}
	for i := 0; i < threads; i++ {
		start := i * batchSize
		if start >= len(pairs) {
			break
		}
		end := int(math.Min(float64(start+batchSize), float64(len(pairs))))
		bucketCollection = append(bucketCollection, make([]bucket, numberOfBuckets))
		bucketSize := maxMorton / uint64(numberOfBuckets)
		wg.Add(1)
		go func(input []mortonPair, threadNumber int) {
			for _, pair := range input {
				index := pair.mortonCode / bucketSize
				bucketCollection[threadNumber][index] = append(bucketCollection[threadNumber][index], pair)
				atomic.AddInt32(&bucketEntries[index], 1)
			}
			wg.Done()
		}(pairs[start:end], i)
	}
	wg.Wait()

	jobs := make(chan mergeJob, threads)
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			for job := range jobs {
				mergeBuckets(job.buckets, job.out)
			}
			wg.Done()
		}()
	}
	start := 0
	for i := 0; i < numberOfBuckets; i++ {
		end := start + int(bucketEntries[i])
		job := mergeJob{
			index: i,
			out:   pairs[start:end],
		}
		for _, buck := range bucketCollection {
			job.buckets = append(job.buckets, buck[i])
		}
		jobs <- job
		start += int(bucketEntries[i])
	}
	close(jobs)
	wg.Wait()
	return pairs
}

// Merges n buckets with the same bucket index
func mergeBuckets(buckets []bucket, out []mortonPair) {
	index := 0
	for _, bucket := range buckets {
		for _, pair := range bucket {
			out[index] = pair
			index++
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].mortonCode < out[j].mortonCode
	})
}
