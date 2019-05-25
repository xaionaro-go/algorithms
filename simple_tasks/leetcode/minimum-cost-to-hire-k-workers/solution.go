package main

import (
	"sort"
)

type candidate struct {
	quality int
	wage    int
}

func (c *candidate) GetEffectiveness() float64 {
	return float64(c.quality) / float64(c.wage)
}

func mincostToHireWorkers(quality []int, wage []int, K int) float64 {
	candidates := make([]*candidate, len(quality))
	for idx := range quality {
		candidates[idx] = &candidate{
			quality[idx],
			wage[idx],
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].GetEffectiveness() > candidates[j].GetEffectiveness()
	})
	theLeastEffectiveCandidate := candidates[K-1]
	theLeastEffectiveEffectiveness := theLeastEffectiveCandidate.GetEffectiveness()

	sum := float64(0)
	for i := 0; i < K; i++ {
		sum += float64(candidates[i].quality) / theLeastEffectiveEffectiveness
	}
	return sum
}
