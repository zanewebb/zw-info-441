package mergesort

import (
	"reflect"
	"testing"
)

// TODO: Add example of cases below and test them
func TestMergeSort(t *testing.T) {
	cases := []struct {
		input          []int
		expectedOutput []int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		},
		{
			[]int{1, 2, 4, 3, 5, 6},
			[]int{1, 2, 3, 4, 5, 6},
		},
		{
			[]int{1, 5, 4, 3, 5, 6},
			[]int{1, 3, 4, 5, 5, 6},
		},
		{
			[]int{1, 1, 1, 1, 1, 1},
			[]int{1, 1, 1, 1, 1, 1},
		},
		{
			[]int{1, 1, 1, 1, 1, 1},
			[]int{1, 1, 1, 1, 1, 1},
		},
		{
			[]int{5, 4, 3, 2, 1, 0},
			[]int{0, 1, 2, 3, 4, 5},
		},
		{
			[]int{-5, 9000, 3, 0, 1, 0},
			[]int{-5, 0, 0, 1, 3, 9000},
		},
		{
			[]int{},
			[]int{},
		},
	}

	for _, c := range cases {
		output := MergeSort(c.input)
		if !reflect.DeepEqual(output, c.expectedOutput) {
			t.Errorf("incorrect output for %v: expected %v but got %v", c.input, c.expectedOutput, output)
		}
	}
}
