package main

import (
	"fmt"
	"math"
)


func main() {
	rate := 0.3725
	for amnt := 0; amnt < 3726; amnt++ {
		get := int(math.Floor(rate * float64(amnt)))
		cost := (float64(get) / rate)
		for cost != math.Trunc(cost) {
			//fmt.Println("nay", cost, get)
			get--
			cost = (float64(get) / rate)
		}
		//fmt.Println("yay", cost, get)
		if get < 1 || cost < 1 { //|| (float64(get)/rate != math.Trunc(float64(get)/rate)){
			fmt.Println("you will get less than 1 dust or need to convert more")
		}
		fmt.Println("for",amnt,"paid", int(cost), "got", get)
	}
}
