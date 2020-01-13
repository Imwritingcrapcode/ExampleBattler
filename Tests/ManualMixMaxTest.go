package main

import (
	. "../Abstract"
	. "../Game"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var DEPTH int
	rand.Seed(time.Now().UTC().UnixNano())
	TWOBOTS := false
	//INIT2 := 9
	//index1 := rand.Intn(len(ReleasedCharacters))
	//index2 := rand.Intn(len(ReleasedCharacters))
	//INIT1 := ReleasedCharacters[index1]
	//INIT2 := ReleasedCharacters[index2]
	//nums := []int{INIT1, INIT2}
	//PLAYER := nums[rand.Intn(len(nums))]
	//for INIT2 == INIT1 {
	//	INIT2 = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	//}
	// Story < E, z < E Ruby > E, Speed > E, Lana > E (1 HP!!), Struc < E
	// Story < z, E > z Ruby > z, Speed > z, Lana > z , Struc > z
	//
	//INIT1 := 1
	INIT2 := 8
	//INIT1 := 9
	//INIT1 := 10
	//INIT1 := 33
	//INIT1 := 51
	INIT1 := 119
	PLAYER := 3000

	if INIT1 == 33 || INIT2 == 33 {
		DEPTH = 5
	} else {
		DEPTH = 6
	}

	/*coin := rand.Intn(2)
	if coin == 0 {
		c := INIT1
		INIT1 = INIT2
		INIT2 = c
	}*/

	var g1, g2 CharInt
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	InitAsNumberBeta(gi1, INIT1)
	InitAsNumberBeta(gi2, INIT2)
	/*input := make(chan string)
	output := make(chan string)*/
	for i := 1; i < 21; i++ {
		if gi1.Number != PLAYER && !gi1.HasEffect(ControlledByStT) ||
			gi1.Number == PLAYER && gi1.HasEffect(ControlledByStT) ||
			TWOBOTS {
			testfood1 := gi1.Copy()
			testfood2 := gi2.Copy()
			prediction, _, moves := MiniMax(testfood1, testfood2, i, DEPTH, true, []int{})
			use := moves[0:2]
			fmt.Println(GetGameState(gi1, gi2, i, gi1.CheckAvailableSkills(i)))
			if gi1.HasEffect(SpedUp) {
				fmt.Println("The bot has used", ToStringStrat([]int{moves[0], moves[1]})+", and predicted", prediction, "\n")
			} else {
				fmt.Println("the bot has used", gi1.Skills[moves[0]].Name+", and predicted", prediction, "\n")
				//fmt.Println("the bot has used", ToStringStrat(moves) + ", and predicted", prediction, "\n")
			}
			TurnApply(gi1, gi2, i, use)

		} else {
			TurnKeyboard(gi1, gi2, i)
		}

		other := gi1
		gi1 = gi2
		gi2 = other
		if !gi1.IsAlive() || !gi2.IsAlive() {
			break

		}
	}
	fmt.Println(g1.(*Girl).Name, g1.(*Girl).CurrHP, g2.(*Girl).Name, g2.(*Girl).CurrHP)

}
