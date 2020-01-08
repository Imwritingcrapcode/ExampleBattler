package main

import (
	. "../Abstract"
	. "../Game"
	"fmt"
)

func main() {
	INIT1 := 51
	INIT2 := 119
	var g1, g2 CharInt
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	InitAsNumber(gi1, INIT1)
	InitAsNumber(gi2, INIT2)
	gi1.Skills[2].Use(gi1, gi2, 1)
	gi2.Skills[1].Use(gi2, gi1, 2)
	gi1.Skills[0].Use(gi1, gi2, 3)
	gi2.Skills[1].Use(gi2, gi1, 4)
	gi1.Skills[0].Use(gi1, gi2, 5)
	gi2.Skills[1].Use(gi2, gi1, 6)
	gi1.Skills[2].Use(gi1, gi2, 7)
	gi2.Skills[0].Use(gi2, gi1, 8)
	gi1.Skills[0].Use(gi1, gi2, 9)
	gi2.Skills[0].Use(gi2, gi1, 10)
	gi1.Skills[0].Use(gi1, gi2, 11)
	gi2.Skills[0].Use(gi2, gi1, 12)
	gi1.Skills[2].Use(gi1, gi2, 13)
	gi2.Skills[3].Use(gi2, gi1, 14)
	gi1.Skills[2].CurrCD = gi1.Skills[2].CD
	fmt.Println("D:", GetGameState(gi1, gi2, 15, gi1.CheckAvailableSkills(15)))
	gi2.DecreaseEffects(gi1, 13)
	gi1.DecreaseEffects(gi2, 14)

	fmt.Println(MiniMax(g1, g2, 15, 6, true, []int{}))

}
