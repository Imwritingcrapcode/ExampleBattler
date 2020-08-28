package main

import (
	. "../Site"
	"fmt"
	"math"
	"strconv"
)

func main() {
	AD := 0
	TOTAL := 1000000
	SP := 0
	LF := 0
	for i := 0; i < TOTAL; i++ {
		//currTime := time.Now().UTC().UnixNano()
		free_girl1, free_girl2 := GenerateStartingPack(int64(i))
		//fmt.Println(i, free_girl1, free_girl2)
		if free_girl1 == 8 || free_girl1 == 9 || free_girl1 == 33 {
			AD += 1
		} else if free_girl1 == 10 || free_girl1 == 51 {
			SP += 1
		} else {
			LF += 1
		}
		if free_girl2 == 8 || free_girl2 == 9 || free_girl2 == 33 {
			AD += 1
		} else if free_girl2 == 10 || free_girl2 == 51 {
			SP += 1
		} else {
			LF += 1
		}
	}
	TOTAL *= 2
	fmt.Println(AD, SP, LF, TOTAL)
	fmt.Println("AD:", strconv.FormatInt(int64(math.Round(float64(float64(AD)/float64(TOTAL)*100))), 10)+"%")
	fmt.Println("SP:", strconv.FormatInt(int64(math.Round(float64(float64(SP)/float64(TOTAL)*100))), 10)+"%")
	fmt.Println("LF:", strconv.FormatInt(int64(math.Round(float64(float64(LF)/float64(TOTAL)*100))), 10)+"%")
}
