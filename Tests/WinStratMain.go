package main

import (
	. "../Abstract"
	. "../Game"
	"fmt"
	"os"
)

func main() {
	INIT1 := 10
	INIT2 := 51
	test := 10

	wins, err := os.OpenFile("C:/Users/~C-o-L/GoglandProjects/Battler/TestFiles/WinStrats.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	losses, err := os.OpenFile("C:/Users/~C-o-L/GoglandProjects/Battler/TestFiles/LoseStrats.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	brokens, err := os.OpenFile("C:/Users/~C-o-L/GoglandProjects/Battler/TestFiles/BrokenStrats.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer wins.Close()
	defer losses.Close()
	defer brokens.Close()
	var g1, g2 CharInt
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	InitAsNumberBeta(gi1, INIT1)
	InitAsNumberBeta(gi2, INIT2)
	LEN := 10
	strat := make([]int, LEN)
	next := make([]int, LEN)
	verdict := -1

	/*Win1.InitWinVert(gi1.Name)
	Win2.InitWinVert(gi2.Name)
	Draw.InitWinVert("Draw!!")
	Tree := Graph{
		false,
		0,
		"[label=\"" + gi1.Name + ": Game Start!" + "\", color=" + GraphColours[Gray] + ", style=filled];\n",
		TimeName(),
		(*Girl)(gi1),
		(*Girl)(gi2),
		make([]*Graph, 1),
	}*/

	if test == gi1.Number {
		wins.WriteString("Testing " + gi1.Name + " against " + gi2.Name + "\n")
		losses.WriteString("Testing " + gi1.Name + " against " + gi2.Name + "\n")
		brokens.WriteString("Testing " + gi1.Name + " against " + gi2.Name + "\n")
	} else {
		wins.WriteString("Testing " + gi2.Name + " against " + gi1.Name + "\n")
		losses.WriteString("Testing " + gi2.Name + " against " + gi1.Name + "\n")
		brokens.WriteString("Testing " + gi2.Name + " against " + gi1.Name + "\n")
	}
	var percReady int
	percReady = 0
	strat = FromStringStrat("qwwqwwer")
	MAXTEST := 1
	//MAXTEST := math.MaxInt64
	for i := 0; i < MAXTEST; i++ {
		if MAXTEST > 1000 && (i-1)%100000 == 0 {
			fmt.Println(percReady, "%")
			percReady += 10
		}
		g1 = new(Girl)
		g2 = new(Girl)
		gi1 = g1.(*Girl)
		gi2 = g2.(*Girl)
		InitAsNumberBeta(gi1, INIT1)
		InitAsNumberBeta(gi2, INIT2)
		enemyStrat := make([]int, 0)
		_, enemyStrat, verdict = TestStrat(gi1, gi2, 1, strat, enemyStrat, test, 1)
		if verdict == 0 {
			losses.WriteString(ToStringStrat(strat) + "; ")
			losses.WriteString("lost to: " + ToStringStrat(enemyStrat) + "\n")
		} else if verdict == 1 || verdict == -1 {
			wins.WriteString(ToStringStrat(strat) + "\n")
			//break
		} else {
			brokens.WriteString(ToStringStrat(strat) + "\n")
		}
		next = NextStrat(strat)
		if next != nil {
			strat = next
		} else {
			break
		}
	}
	wins.WriteString("\n")
	losses.WriteString("\n")
	brokens.WriteString("\n")
	fmt.Println(Combinations)
	/*graph, err := os.OpenFile("C:/Users/~C-o-L/GoglandProjects/Battler/TestFiles/graph.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	graph.WriteString(Tree.ToString())*/

	fmt.Println("done")

}
