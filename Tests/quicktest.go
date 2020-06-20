package main

import (
	"math"
	"fmt"
)

func main() {
	played := 16
	won := 3
	var winrate float64
	if won == 0 {
		winrate = 0
	} else {
		winrate = math.Floor(float64(won)/float64(played)*100)
	}
	fmt.Println(winrate)
	switch {
	case winrate > 50:
		fmt.Println(6)
	case winrate >= 40:
		fmt.Println(5)
	case winrate >= 30:
		fmt.Println(4)
	default:
		fmt.Println(3)
	}




}