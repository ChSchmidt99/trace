package pt

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

// TODO: Merge file with LBVH

type bucket []mortonPair

type mergeJob struct {
	index   int
	buckets []bucket
	out     []mortonPair
}

// Parallel bucket sort
func sortMortonPairs(pairs []mortonPair, numberOfBuckets int, threads int) {
	bucketCollection, bucketFill := fillBuckets(pairs, numberOfBuckets, threads)
	merge(pairs, bucketFill, bucketCollection, numberOfBuckets, threads)
}

// Inserts morton pairs into the specified number of buckets
// Each thread uses a separate slice of buckets to avoid the need for synchronized access
// Return:
// buckets: [threads][numberOfBuckets]bucket => one slice of buckets for each thread
// bucketFill: holds how many pairs have been inserted into the corresponding bucket
func fillBuckets(pairs []mortonPair, numberOfBuckets int, threads int) (buckets [][]bucket, bucketFill []int32) {
	batchSize := int(math.Ceil(float64(len(pairs)) / float64(threads)))
	bucketCollection := make([][]bucket, 0, threads)
	bucketEntries := make([]int32, numberOfBuckets)
	wg := sync.WaitGroup{}

	// Each thread inserts an equal amount of pairs into its seperate slice of buckets
	for i := 0; i < threads; i++ {
		start := i * batchSize
		if start >= len(pairs) {
			break
		}
		end := int(math.Min(float64(start+batchSize), float64(len(pairs))))
		bucketCollection = append(bucketCollection, make([]bucket, numberOfBuckets))
		bucketSize := MAX_MORTON_CODE / uint64(numberOfBuckets)
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
	return bucketCollection, bucketEntries
}

func merge(out []mortonPair, bucketEntries []int32, bucketCollection [][]bucket, numberOfBuckets int, threads int) {
	// Start workers, each worker inserts pairs into the given interval of the out slice and sorts it
	jobs := make(chan mergeJob, threads)
	wg := sync.WaitGroup{}
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			for job := range jobs {
				mergeBuckets(job.buckets, job.out)
			}
			wg.Done()
		}()
	}

	// Feed jobs to workers,
	// Bucket fills are used to determine the corresponding interval in the output slice
	// This method is used to avoid allocating a output slice as this would be quite expensive
	start := 0
	for i := 0; i < numberOfBuckets; i++ {
		end := start + int(bucketEntries[i])
		job := mergeJob{
			index: i,
			out:   out[start:end],
		}
		for _, buck := range bucketCollection {
			job.buckets = append(job.buckets, buck[i])
		}
		jobs <- job
		start += int(bucketEntries[i])
	}
	close(jobs)
	wg.Wait()
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
