package main

import (
	"fmt"
	"math"
)

type shape interface {
	calcArea() float64
}

type circle struct {
	radius int
}

type square struct {
	side int
}

func (c circle) calcArea() float64 {
	return float64(c.radius*c.radius) * math.Pi
}

func (s square) calcArea() float64 {
	return float64(s.side * s.side)
}

func main() {
	c := circle{
		radius: 5,
	}

	s := square{
		side: 3,
	}

	info(c)
	info(s)
}

func info(s shape) {
	fmt.Println(s.calcArea())
}
