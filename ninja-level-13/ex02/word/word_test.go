package word_test

import (
	"fmt"
	"testing"

	"github.com/preslavmihaylov/learn-golang/ninja-level-13/ex02/word"
)

func TestCount(t *testing.T) {
	res := word.Count("Here's Johnny!")
	if res != 2 {
		t.Errorf("Expected 2, got %v", res)
	}
}

func ExampleCount() {
	fmt.Println(word.Count("one two three"))
	// Output:
	// 3
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		word.Count("Here's Johnny!")
	}
}

func TestUseCount(t *testing.T) {
	res := word.UseCount("one two three two three")
	if res["one"] != 1 {
		t.Errorf("Expected map[\"one\"] = 1, got map[\"one\"] = %v", res["one"])
	}

	if res["two"] != 2 {
		t.Errorf("Expected map[\"two\"] = 2, got map[\"two\"] = %v", res["two"])
	}

	if res["three"] != 2 {
		t.Errorf("Expected map[\"three\"] = 2, got map[\"three\"] = %v", res["three"])
	}
}

func BenchmarkUseCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		word.UseCount("Here's Johnny!")
	}
}
