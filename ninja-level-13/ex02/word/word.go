// Package word provides functions for handling word queries in a string
package word

import "strings"

// UseCount returns a map representing the word occurences for each word in the provided string
func UseCount(s string) map[string]int {
	xs := strings.Fields(s)
	m := make(map[string]int)
	for _, v := range xs {
		m[v]++
	}
	return m
}

// Count returns the total word count of the provided string
func Count(s string) int {
	return len(strings.Fields(s))
}
