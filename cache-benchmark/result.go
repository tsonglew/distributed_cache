package main

import (
	"time"
)

type result struct {
	getCount    int
	missCount   int
	setCount    int
	statBuckets []statistic
}

func (r *result) addResult(src *result) {
	for b, s := range src.statBuckets {
		r.addStatistic(b, s)
	}
	r.getCount += src.getCount
	r.setCount += src.setCount
	r.missCount += src.missCount
}

func (r *result) addDuration(d time.Duration, typ string) {
	bucket := int(d / time.Millisecond)
	r.addStatistic(bucket, statistic{1, d})
	switch typ {
	case "get":
		r.getCount++
	case "set":
		r.setCount++
	default:
		r.missCount++
	}
}

func (r *result) addStatistic(bucketIdx int, stat statistic) {
	// for scale
	if bucketIdx > len(r.statBuckets)-1 {
		newBuckets := make([]statistic, bucketIdx+1)
		copy(newBuckets, r.statBuckets)
		r.statBuckets = newBuckets
	}
	r.statBuckets[bucketIdx].count += stat.count
	r.statBuckets[bucketIdx].time += stat.time
}
