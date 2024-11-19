package stats

import "markov_generator/domain"

type KigoAveragePos float32
type NumHaiku int

type KigoInfo struct {
	Kigo       domain.Kigo
	NumHaiku   int
	AveragePos KigoAveragePos
}

type KigoStat struct {
	TotalAverage KigoAveragePos
	Individual   map[domain.Season]([]*KigoInfo)
}
