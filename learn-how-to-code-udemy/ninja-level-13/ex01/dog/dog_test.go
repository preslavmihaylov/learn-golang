package dog_test

import (
	"fmt"
	"testing"

	"github.com/preslavmihaylov/learn-golang/ninja-level-13/ex01/dog"
)

var yearsTests = [][]int{
	[]int{10, 70},
	[]int{5, 35},
	[]int{1, 7},
}

func TestYears(t *testing.T) {
	for _, v := range yearsTests {
		input := v[0]
		output := v[1]
		dogYears := dog.Years(input)
		if dogYears != output {
			t.Errorf("Expected %v, Got %v", output, dogYears)
		}
	}
}

func ExampleYears() {
	fmt.Println(dog.Years(10))
	// Output:
	// 70
}

func BenchmarkYears(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dog.Years(10)
	}
}

var yearsTwoTests = [][]int{
	[]int{10, 70},
	[]int{5, 35},
	[]int{1, 7},
}

func TestYearsTwo(t *testing.T) {
	for _, v := range yearsTests {
		input := v[0]
		output := v[1]
		dogYears := dog.YearsTwo(input)
		if dogYears != output {
			t.Errorf("Expected %v, Got %v", output, dogYears)
		}
	}
}

func ExampleYearsTwo() {
	fmt.Println(dog.YearsTwo(10))
	// Output:
	// 70
}

func BenchmarkYearsTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dog.YearsTwo(10)
	}
}
