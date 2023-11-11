package gamilton

import "slices"

const NEAR_INFINITE_NUMBER = 99999

func HungryAlgorythm(matrix [][]float64) (path []int) {
	path = make([]int, len(matrix))
	path[0] = 0
	path[len(path)-1] = 1
	checkingVertex := 0
	for i := 0; i < len(matrix)-1; i++ {
		min := float64(NEAR_INFINITE_NUMBER)
		minIdx := checkingVertex
		for k, v := range matrix[checkingVertex] {
			if !slices.Contains(path, k) {
				if v < min {
					min = v
					minIdx = k
				}
			}
		}
		path[i+1] = minIdx
		checkingVertex = i
	}

	return
}
