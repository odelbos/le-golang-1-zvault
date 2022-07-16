package main

import (
	"math"
)

const (
	BLOCK_SIZE = 4096
	KO         = 1024
	MO         = 1024 * 1024
	GO         = 1024 * 1024 * 1024
	TO         = 1024 * 1024 * 1024
)

func NbGroups(size int64) float64 {
	if size < 100*KO {
		return 10
	}
	segments := getStructure()
	for _, segment := range segments {
		x2 := segment[2]
		if size < x2 {
			x1, y1, y2 := segment[0], segment[1], segment[3]
			m := (float64(y2) - float64(y1)) / (float64(x2) - float64(x1))
			b := float64(y1) - m*float64(x1)
			return math.Ceil(m*float64(size) + b)
		}
	}
	return 20000
}

func NbBlocksPerGroup(size int64) (int, int) {
	nbGroups := NbGroups(size)
	gs := float64(size) / nbGroups
	nbBlocks := math.Ceil(gs / BLOCK_SIZE)
	return int(nbBlocks), BLOCK_SIZE
}

//
// Helper functions
//
func getStructure() [][]int64 {
	return [][]int64{
		// x1, y1, x2, y2
		{100 * KO, 10, MO, 20},
		{MO, 20, 10 * MO, 30},
		{10 * MO, 30, 100 * MO, 40},
		{100 * MO, 40, 500 * MO, 50},
		{500 * MO, 50, GO, 70},
		{GO, 70, 10 * GO, 120},
		{10 * GO, 120, 100 * GO, 250},
		{100 * GO, 250, 500 * GO, 500},
		{500 * GO, 500, TO, 1000},
		{TO, 1000, 10 * TO, 2000},
		{10 * TO, 2000, 50 * TO, 4000},
		{50 * TO, 4000, 100 * TO, 8000},
		{100 * TO, 8000, 150 * TO, 12000},
		{150 * TO, 12000, 250 * TO, 20000}}
}
