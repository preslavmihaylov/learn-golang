package mymath_test

import (
	"fmt"
	"testing"

	"github.com/preslavmihaylov/learn-golang/ninja-level-13/ex03/mymath"
)

type testCenteredAvgData struct {
	input  []int
	output float64
}

var testsCenteredAvg = []testCenteredAvgData{
	testCenteredAvgData{[]int{1, 4, 6, 8, 100}, 6},
	testCenteredAvgData{[]int{0, 8, 10, 1000}, 9},
	testCenteredAvgData{[]int{9000, 4, 10, 8, 6, 12}, 9},
	testCenteredAvgData{[]int{123, 744, 140, 200}, 170},
}

func TestCenteredAvg(t *testing.T) {
	for _, v := range testsCenteredAvg {
		res := mymath.CenteredAvg(v.input)
		if res != v.output {
			t.Errorf("Expected %v, got %v", v.output, res)
		}
	}
}

func ExampleCenteredAvg() {
	fmt.Println(mymath.CenteredAvg([]int{1, 2, 3, 4, 5, 6}))
	// Output:
	// 3.5
}

func BenchmarkCenteredAvg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mymath.CenteredAvg(testsCenteredAvg[i%len(testsCenteredAvg)].input)
	}
}
