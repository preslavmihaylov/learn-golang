package main

import "fmt"

type vehicle struct {
	doorsCnt int
	color    string
}

type truck struct {
	vehicle
	fourWheel bool
}

type sedan struct {
	vehicle
	luxury bool
}

func main() {
	myTruck := truck{
		vehicle: vehicle{
			doorsCnt: 4,
			color:    "Green",
		},
		fourWheel: true,
	}

	mySedan := sedan{
		vehicle: vehicle{
			doorsCnt: 2,
			color:    "Red",
		},
		luxury: false,
	}

	fmt.Println(myTruck)
	fmt.Println(mySedan)
	fmt.Println("My truck's color: ", myTruck.color)
	fmt.Println("My sedan's color: ", mySedan.color)
}
