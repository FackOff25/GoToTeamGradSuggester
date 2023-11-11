package gamilton_test

import (
	"fmt"
	"testing"

	"github.com/FackOff25/GoToTeamGradSuggester/pkg/gamilton"
)

func TestHungryAlgorythm(t *testing.T) {
	places := [][]float64{
		{9, 1, 3, 2},
		{1, 9, 1, 2},
		{3, 1, 999, 3},
		{2, 2, 3, 999},
	}

	path := gamilton.HungryAlgorythm(places)

	fmt.Printf("%v", path)
}
